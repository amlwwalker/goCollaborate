package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"io/ioutil"
)

/*
Every time a new client connects, we need to read the current contents of the file
*/

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	ws *websocket.Conn
	send chan []byte // Channel storing outgoing messages
	id string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

//a new connection is made
func serveWs(w http.ResponseWriter, r *http.Request) {
	identity := r.URL.Query().Get("id")
	filename := r.URL.Query().Get("filename")
	log.Println("filename: ", filename)
	// userRepo.identity = identity
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)

	c := &client{
		send: make(chan []byte, maxMessageSize),
		ws: ws,
		id: identity,
	}
	log.Println("Created client: " + c.id)
	// read the content of the file when a user connects. They need the latest version from master
 	dat, err := ioutil.ReadFile(location + "/"+filename)
    check(err)
    log.Println(string(dat))
  	h.content = string(dat)

	h.register <- c

	go c.writePump()
	c.readPump(filename)
}

func (c *client) readPump(filename string) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		check(err)
		log.Println(c.id + " writing to file: " + filename)
		err = ioutil.WriteFile(location + "/"+filename, message, 0644)
		// userRepo.stageChanges(filename)
		// t := time.Now()
		// userRepo.commitChanges(c.id + " commited to " + filename + " at " + t.Format("20060102150405"))
		check(err)
		h.broadcast <- string(message)

	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		log.Println("WTF: " + c.id)
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}