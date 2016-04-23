// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  //"encoding/json"
  "github.com/gorilla/websocket"
  "net/http"
  //"log"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	s *switcher
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
	h    *hub
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

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

  room := get_room(wsh.s, r.Header.Get("X-Room"))
	c := &connection{send: make(chan []byte, 256), ws: ws, h: room}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()

	go c.writer()
	c.reader()
}
