package models

import (
	"encoding/json"
)

type Hub struct {
	ClientList  map[*Client]bool
	ChBoardcast chan []byte
	ChRegister  chan *Client
	ChLogout    chan *Client
}

func NewHub() *Hub {
	return &Hub{
		ClientList:  make(map[*Client]bool),
		ChBoardcast: make(chan []byte),
		ChRegister:  make(chan *Client),
		ChLogout:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.ChRegister:
			h.ClientList[client] = true
			client.Ip = client.Conn.RemoteAddr().String()
			client.ChatData.Type = "handshake"
			client.ChatData.UserList = User_list
			RegisterMessage, _ := json.Marshal(client.ChatData)
			client.ChSend <- RegisterMessage

		case client := <-h.ChLogout:
			if _, ok := h.ClientList[client]; ok {
				delete(h.ClientList, client)
				close(client.ChSend)
			}

		case chatdata := <-h.ChBoardcast:
			for client := range h.ClientList {
				select {
				case client.ChSend <- chatdata:
				default:
					delete(h.ClientList, client)
					close(client.ChSend)
				}
			}
		}
	}
}
