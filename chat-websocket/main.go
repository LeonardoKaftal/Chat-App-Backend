package main

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"sync"
)

type Message struct {
	payload         []byte
	destinationRoom string
	sender          *websocket.Conn
}

type ConnectedUser struct {
	conn     *websocket.Conn
	username string
	roomName string
}

func NewUser(conn *websocket.Conn, username string, roomName string) *ConnectedUser {
	return &ConnectedUser{
		conn:     conn,
		username: username,
		roomName: roomName,
	}
}

type Server struct {
	globalRoomConns   map[*ConnectedUser]bool
	privateRooms      map[string]map[*ConnectedUser]bool
	mut               sync.Mutex
	msgChannel        chan Message
	ConnectedUsername mapset.Set[*ConnectedUser]
}

func NewServer() *Server {
	return &Server{
		globalRoomConns:   make(map[*ConnectedUser]bool),
		privateRooms:      make(map[string]map[*ConnectedUser]bool),
		msgChannel:        make(chan Message),
		ConnectedUsername: mapset.NewSet(&ConnectedUser{}),
	}
}

func (s *Server) addConnection(user *ConnectedUser) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.ConnectedUsername.Add(user)
	s.broadCastConnectedUser(user.username)
	if user.roomName != "" {
		if s.privateRooms[user.roomName] == nil {
			s.privateRooms[user.roomName] = make(map[*ConnectedUser]bool)
		}
		s.privateRooms[user.roomName][user] = true
	} else {
		s.globalRoomConns[user] = true
	}
}

func (s *Server) closeConnection(user *ConnectedUser) {
	s.mut.Lock()
	fmt.Printf("User %s disconnected\n", user.username)
	defer s.mut.Unlock()
	s.ConnectedUsername.Remove(user)
	delete(s.globalRoomConns, user)
	delete(s.privateRooms[user.roomName], user)
	if len(s.privateRooms) == 0 {
		delete(s.privateRooms, user.roomName)
	}
	err := user.conn.Close()
	if err != nil {
		fmt.Println("Error attempting to close the connection for user " + user.username + ", The error is the following " + err.Error())
		return
	}
}

func (s *Server) readLoop(user *ConnectedUser) {
	defer s.closeConnection(user)
	buff := make([]byte, 2048)
	for {
		n, err := user.conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error trying to read the message ", err.Error())
		}
		s.msgChannel <- Message{
			payload:         buff[:n],
			destinationRoom: user.roomName,
			sender:          user.conn,
		}
	}
}

func (s *Server) messageHandler() {
	for msg := range s.msgChannel {
		if msg.destinationRoom == "" {
			for user := range s.globalRoomConns {
				if user.conn != msg.sender {
					if _, err := user.conn.Write(msg.payload); err != nil {
						println("Error broadcasting the message in the global room " + err.Error())
					}
				}
			}
		} else {
			for user := range s.privateRooms[msg.destinationRoom] {
				if user.conn != msg.sender {
					if _, err := user.conn.Write(msg.payload); err != nil {
						println("Error broadcasting the message in the global room " + err.Error())
					}
				}
			}
		}
	}
}

func (s *Server) broadCastConnectedUser(username string) {
	usernames := make([]string, 0)

	for user := range s.ConnectedUsername.Iter() {
		usernames = append(usernames, user.username)
	}

	data, err := json.Marshal(usernames)
	if err != nil {
		fmt.Println("Error encoding usernames: ", err)
		return
	}

	for user := range s.ConnectedUsername.Iter() {
		if _, err := user.conn.Write(data); err != nil {
			fmt.Println("Error broadcasting usernames: ", err)
		}
	}
}

var server = NewServer()

func WebsocketHandler(conn *websocket.Conn) {
	vars := mux.Vars(conn.Request())
	username := vars["username"]
	roomName := vars["room"]
	user := NewUser(conn, username, roomName)
	fmt.Printf("Incoming connection from %s with username: %s\n", conn.RemoteAddr(), username)
	server.addConnection(user)
	server.broadCastConnectedUser(username)
	server.readLoop(user)
}

func main() {
	server.ConnectedUsername.Clear()
	r := mux.NewRouter()
	r.Handle("/ws/{username}", websocket.Handler(WebsocketHandler))
	r.Handle("/ws/{username}/{room}", websocket.Handler(WebsocketHandler))
	http.Handle("/", r)

	go server.messageHandler()
	log.Fatal(http.ListenAndServe(":8090", nil))
}
