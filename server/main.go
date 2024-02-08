package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clientConnections = make(map[string]*websocket.Conn)

type Message struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	user_id := r.URL.Query().Get("user_id")
	log.Println("user id:", user_id)
	clientConnections[user_id] = conn

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		var msg Message
		err = json.Unmarshal(b, &msg)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if wc, ok := clientConnections[msg.To]; ok {
			wc.WriteMessage(websocket.TextMessage, b)
		} else {
			log.Println("user id not connected: ", msg.To)
		}

	}
}

func Run() {
	log.Println("Server is running")

	http.HandleFunc("/ws", websocketHandler)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
