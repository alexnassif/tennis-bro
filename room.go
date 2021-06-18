package main

import (
	"fmt"
	"github.com/google/uuid"
)
const welcomeMessage = "%s joined the room"
type Room struct {

	ID 			uuid.UUID `json:"id"`
	Name 		string `json:"name"`
	Private 	bool `json:"private"`
	clients 	map[*Client]bool
	register 	chan *Client
	unregister 	chan *Client
	broadcast 	chan *Message
}

//create a new Room
func NewRoom(name string, private bool) *Room {
	fmt.Print(name)
	return &Room{
		ID: uuid.New(),
		Name: name,
		Private: private,
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan *Message),
	}
}

//run room, accept various requests
func (room *Room) RunRoom(){
	for{
		select{
		case client := <-room.register:
			room.registerClientInRoom(client)
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client){
	if !room.Private{
		room.notifyClientJoined(client)
	}
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client){
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
    return room.Private
}
