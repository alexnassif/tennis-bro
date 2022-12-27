package main

import (
	"fmt"

	"github.com/alexnassif/tennis-bro/Models"
	"github.com/google/uuid"

	"encoding/json"
	"log"

	"github.com/alexnassif/tennis-bro/Config"
)

//import "golang.org/x/text/message"

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
	users      []Models.OnlineUser
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	wsServer := &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
	var users []Models.OnlineClient

	Models.GetAllOnlineUsers(&users)

	wsServer.users = make([]Models.OnlineUser, 0)
	for _, val := range users {
		wsServer.users = append(wsServer.users, &val)
	}
	
	return wsServer
}

func (server *WsServer) Run() {
	go server.listenPubSubChannel()
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	//add user to online table
	onlineUser := Models.OnlineClient{ID: client.User.GetId(), User: client.User}
	Models.AddOnlineClient(&onlineUser)
	//server.notifyClientJoined(client)
	server.users = append(server.users, &onlineUser)
	server.publishClientJoined(client)

	server.listOnlineClients(client)
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)
		onlineUser := Models.OnlineClient{ID: client.User.GetId(), User: client.User}
		Models.RemoveOnlineUser(&onlineUser)

		server.publishClientLeft(client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	fmt.Println("finding room" + name)
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			fmt.Println("break")
			break
		}
	}

	return foundRoom
}

func (server *WsServer) runRoomFromRepository(name string) *Room {
	var room Models.Room
	var newRoom *Room
	err := Models.FindRoomByName(&room, name)
	if err == nil {
		newRoom = NewRoom(room.GetName(), room.GetPrivate())
		newRoom.ID, _ = uuid.Parse(fmt.Sprint(room.GetId()))
		go newRoom.RunRoom()
		server.rooms[newRoom] = true

	}
	return newRoom
}

func (server *WsServer) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)

	newRoom := Models.Room{Name: room.GetName(), Private: room.GetPrivate()}
	Models.AddRoom(&newRoom)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}

func (server *WsServer) createPrivateRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	go room.RunRoom()
	server.rooms[room] = true
	fmt.Println("created room" + name)
	return room
}

func (server *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}
	server.broadcastToClients(message.encode())
}

func (server *WsServer) notifyClientLeft(client *Client) {
	message := &Message{

		Action: UserLeftAction,
		Sender: client,
	}
	server.broadcastToClients(message.encode())
}

func (server *WsServer) listOnlineClients(client *Client) {

	for _, user := range server.users {
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.send <- message.encode()
	}
}

func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

// Publish userJoined message in pub/sub
func (server *WsServer) publishClientJoined(client *Client) {

	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	if err := Config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Publish userleft message in pub/sub
func (server *WsServer) publishClientLeft(client *Client) {

	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	if err := Config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Listen to pub/sub general channels
func (server *WsServer) listenPubSubChannel() {

	pubsub := Config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()
	for msg := range ch {

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		switch message.Action {
		case UserJoinedAction:
			server.handleUserJoined(message)
		case UserLeftAction:
			server.handleUserLeft(message)
		case JoinRoomPrivateAction:
			server.handleUserJoinPrivate(message)
		}
	}
}

func (server *WsServer) handleUserJoined(message Message) {
	// Add the user to the slice
	server.broadcastToClients(message.encode())
}

func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		if user.GetId() == message.Sender.GetId() {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}
	server.broadcastToClients(message.encode())
}

func (server *WsServer) handleUserJoinPrivate(message Message) {
	// Find client for given user, if found add the user to the room.
	targetClient := server.findClientByID(message.Message)
	if targetClient != nil {
		targetClient.joinRoom(message.Target.GetName(), message.Sender)
	}
}

// Add the findUserByID method used by client.go
func (server *WsServer) findUserByID(ID string) Models.OnlineUser {
	var foundUser Models.OnlineUser
	for _, client := range server.users {
		//fmt.Println(client)
		if client.GetId() == ID {
			foundUser = client
			break
		}
	}

	return foundUser
}
