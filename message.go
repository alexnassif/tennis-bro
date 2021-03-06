package main

import (
	"encoding/json"
	"log"

	"github.com/alexnassif/tennis-bro/Models"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const PrivateMessage = "private-message"

type Message struct {
	Action   string            `json:"action"`
	Message  string            `json:"message"`
	Target   *Room             `json:"target"`
	Sender   Models.OnlineUser `json:"sender"`
	Receiver int            `json:"receiver"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
