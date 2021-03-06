package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexnassif/tennis-bro/Config"
	"github.com/alexnassif/tennis-bro/Models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (

	//Max wait time when writing message to peer
	writeWait = 10 * time.Second

	//Max time until mext pong from peer
	pongWait = 60 * time.Second

	//Send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10

	//Maximum message size allowed from peer
	maxMessageSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client represents the websocket client at the server
type Client struct {
	//websocket connection
	ID       uuid.UUID `json:"id"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	Name     string `json:"name"`
	User     Models.User
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string, user Models.User) *Client {
	return &Client{
		ID:       uuid.New(),
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
		Name:     name,
		User:     user,
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["id"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	var user Models.User
	err := Models.GetUserByID(&user, strings.Join(name, ""))
	if err != nil {
		fmt.Println("no user with that id")
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, user.UserName, user)
	client.joinPrivateRoom(user.GetId(), client)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client

	fmt.Println("New Client joined the hub!")
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = client
	//fmt.Print(string(jsonMessage))

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomID := message.Target.GetId()
		// Use the ChatServer method to find the room, and if found, broadcast!
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}
	// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	case PrivateMessage:
		client.handlePrivateMessage(message)
	}

}

func (client *Client) handlePrivateMessage(message Message) {
	var user Models.User
	err := Models.GetUserByID(&user, fmt.Sprint(message.Receiver))
	if err != nil {
		fmt.Println("no user with that id")
	} else {
		fmt.Println(user.UserName)
	}
	var privateRoom Models.Room
	result := Config.DB.Where("user1_id = ? AND user2_id = ?", user.ID, client.User.ID).Or("user1_id = ? AND user2_id = ?", client.User.ID, user.ID).First(&privateRoom)

	if result.RowsAffected == 0 {

		privateRoom = Models.Room{Name: user.UserName + client.User.UserName, User1: client.User, User2: user, Private: true}
		err := Config.DB.Create(&privateRoom)

		if err.Error != nil {
			fmt.Println(err.Error)
		} else {
			fmt.Println(privateRoom.User1.UserName)
		}
	} else {
		fmt.Println(privateRoom.CreatedAt)
	}

	newMessage := Models.Message{Sender: client.User, Recipient: user, Body: message.Message}
	Config.DB.Create(&newMessage)

	//find room by id
	room := client.wsServer.findRoomByName(user.GetId())
	room.broadcast <- &Message{
		Action:  PrivateMessage,
		Message: newMessage.Body,
		Target:  room,
		Sender:  client,
	}
	
}

func (client *Client) joinPrivateRoom(roomName string, sender Models.OnlineUser) *Room {

	
	room := client.wsServer.createPrivateRoom(roomName, sender != nil)
	

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
	}
	return room
}

// Send out invite message over pub/sub in the general channel.
func (client *Client) sendTargetUserMessage(target Models.OnlineUser, room *Room) {
	inviteMessage := &Message{
		Action:  PrivateMessage,
		Message: target.GetId(),
		Target:  room,
		Sender:  client,
	}

	if err := Config.Redis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Refactored method
// Use new joinRoom method
func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	client.joinRoom(roomName, nil)
}

// Refactored method
// Added nil check
func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByID(message.Message)
	if room == nil {
		return
	}
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

// New method
// When joining a private room we will combine the IDs of the users
// Then we will bothe join the client and the target.
func (client *Client) handleJoinRoomPrivateMessage(message Message) {

	target := client.wsServer.findUserByID(message.Message)
	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	roomName := message.Message + client.ID.String()

	joinedRoom := client.joinRoom(roomName, target)
	//target.joinRoom(roomName, client)
	if joinedRoom != nil {
		client.inviteTargetUser(target, joinedRoom)
	}
}

// Send out invite message over pub/sub in the general channel.
func (client *Client) inviteTargetUser(target Models.OnlineUser, room *Room) {
	inviteMessage := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetId(),
		Target:  room,
		Sender:  client,
	}

	if err := Config.Redis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// New method
// Joining a room both for public and private roooms
// When joiing a private room a sender is passed as the opposing party
func (client *Client) joinRoom(roomName string, sender Models.OnlineUser) *Room {

	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room, sender)
	}
	return room
}

// New method
// Check if the client is not yet in the room
func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}

// New method
// Notify the client of the new room he/she joined
func (client *Client) notifyRoomJoined(room *Room, sender Models.OnlineUser) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) GetId() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) GetUser() Models.User {
	return client.User
}
