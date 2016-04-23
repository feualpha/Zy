// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "encoding/base64"
  "github.com/gorilla/websocket"
  "net/http"
  "strings"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	s *switcher
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
	h    *hub
  name string
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.h.broadcast <- &mesg{body: message, id: c}
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func get_room(s *switcher,id string) *hub{
  y := newSwitchAgent(id)
  var room *hub

  s.join <- *y
  select {
  case h := <- y.room:
    room = h
    break
  }

  return room
}

func get_username(auth_header string) string {
  encoded :=  strings.Split(auth_header, " ")
  decoded,_ := base64.StdEncoding.DecodeString(encoded[1])
  n := len(decoded)
  userpass := string(decoded[:n])
  username := strings.Split(userpass, ":")

  return username[0]
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
  name := get_username(r.Header.Get("Authorization"))
  room := get_room(wsh.s, r.Header.Get("X-Room"))

	c := &connection{send: make(chan []byte, 256), ws: ws, h: room, name: name}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()

	go c.writer()
	c.reader()
}
