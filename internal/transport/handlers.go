package transport

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//TODO: add disconnect action
//TODO: deal with reconnection

var views = jet.NewSet(jet.NewOSFileSystemLoader("./ui/html"), jet.InDevelopmentMode())
var wsChannel = make(chan WsRequest)
var users = make(map[WsConnection]string)

func DisplayIndexPage(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "index.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

var connection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsConnection struct {
	*websocket.Conn
}

type WsRequest struct {
	Action     string       `json:"action"`
	Payload    string       `json:"payload"`
	Connection WsConnection `json:"connection"`
}

type WsResponse struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

func Connect(w http.ResponseWriter, r *http.Request) {
	ws, err := connection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connection := WsConnection{ws}
	users[connection] = ""
	InitialBroadcast(connection)

	go Listen(&connection)
}

func Listen(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Restarting connection:", r)
		}
	}()

	var request WsRequest

	for {
		conn.ReadJSON(&request)
		request.Connection = *conn
		wsChannel <- request
	}

}

func ProcessMessages() {
	var response WsResponse

	for {
		e := <-wsChannel
		switch e.Action {
		case "username":
			if !usernameExists(e.Payload) {
				users[e.Connection] = e.Payload
				response.Action = "update_users"
				response.Payload = GetUserList()
				BroadcastToAll(response)
			} else {
				response.Action = "username_error"
				response.Payload = "Username already exists"
				BroadcastToOne(response, e.Connection)
			}
		case "message":
			response.Action = "message"
			response.Payload = fmt.Sprintf("<strong>%s</strong>: %s", users[e.Connection], e.Payload)
			BroadcastToAll(response)
		}
	}
}

func BroadcastToAll(response WsResponse) {
	for user := range users {
		err := user.WriteJSON(response)
		if err != nil {
			_ = user.Close()
			delete(users, user)
		}
	}
}

func BroadcastToOne(response WsResponse, connection WsConnection) {
	err := connection.WriteJSON(response)
	if err != nil {
		_ = connection.Close()
		delete(users, connection)
	}
}

func InitialBroadcast(connection WsConnection) {
	var response WsResponse
	response.Action = "update_users"
	response.Payload = GetUserList()
	BroadcastToOne(response, connection)
}
