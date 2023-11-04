package main

import (
	"fmt"
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

type Server struct {
	globalRoomConns map[*websocket.Conn]bool
	privateRooms    map[string]map[*websocket.Conn]bool
	mut             sync.Mutex
	msgChannel      chan Message
}

func NewServer() *Server {
	return &Server{
		globalRoomConns: make(map[*websocket.Conn]bool),
		privateRooms:    make(map[string]map[*websocket.Conn]bool),
		msgChannel:      make(chan Message),
	}
}

func (s *Server) addConnection(roomName string, ws *websocket.Conn) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if roomName != "" {
		if s.privateRooms[roomName] == nil {
			s.privateRooms[roomName] = make(map[*websocket.Conn]bool)
		}
		s.privateRooms[roomName][ws] = true
	} else {
		s.globalRoomConns[ws] = true
	}
}

func (s *Server) closeConnection(ws *websocket.Conn, roomName string) {
	s.mut.Lock()
	fmt.Printf("User %s disconnected\n", ws.RemoteAddr())
	defer s.mut.Unlock()
	delete(s.globalRoomConns, ws)
	delete(s.privateRooms[roomName], ws)
	if len(s.privateRooms) == 0 {
		delete(s.privateRooms, roomName)
	}
	err := ws.Close()
	if err != nil {
		fmt.Println("Error attempting to close the connection")
		return
	}
}

func (s *Server) readLoop(ws *websocket.Conn, roomName string) {
	defer s.closeConnection(ws, roomName)
	buff := make([]byte, 2048)
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error trying to read the message ", err.Error())
		}
		s.msgChannel <- Message{
			payload:         buff[:n],
			destinationRoom: roomName,
			sender:          ws,
		}
	}
}

func (s *Server) messageHandler() {
	for msg := range s.msgChannel {
		if msg.destinationRoom == "" {
			for ws := range s.globalRoomConns {
				if ws != msg.sender {
					if _, err := ws.Write(msg.payload); err != nil {
						println("Error broadcasting the message in the global room " + err.Error())
					}
				}
			}
		} else {
			for ws := range s.privateRooms[msg.destinationRoom] {
				if ws != msg.sender {
					if _, err := ws.Write(msg.payload); err != nil {
						println("Error broadcasting the message in the global room " + err.Error())
					}
				}
			}
		}
	}
}

var server = NewServer()

func WebsocketHandler(conn *websocket.Conn) {
	fmt.Println("Incoming connection from ", conn.RemoteAddr().String())
	vars := mux.Vars(conn.Request())
	roomName := vars["room"]
	server.addConnection(roomName, conn)
	server.readLoop(conn, roomName)
}

func main() {
	r := mux.NewRouter()
	r.Handle("/ws/{room}", websocket.Handler(WebsocketHandler))
	r.Handle("/ws", websocket.Handler(WebsocketHandler))
	http.Handle("/", r)

	go server.messageHandler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
