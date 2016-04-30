// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type broadcast_message struct {
	body []byte
	id *connection
}

//this type is contract between client and server
//server message_receive --> client message_send
type message_receive struct {
	Body []byte
}
//server message_send --> client message_client
type message_send struct {
	Sender []byte
	Body []byte
}
//////////////////////////////////////////////////
type hub struct {
	connections map[*connection]bool
	broadcast   chan *broadcast_message
	register    chan *connection
	unregister  chan *connection
	live        chan bool
	special     bool
}

func newHub(special bool) *hub {
	return &hub{
		broadcast:   make(chan *broadcast_message),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
		live: make(chan bool),
		special: special,
	}
}

func unregister(h *hub, c *connection){
	if _, ok := h.connections[c]; ok {
		delete(h.connections, c)
		close(c.send)
	}
}

func will_self_destroy(count int, special bool) bool {
	return !((count==0) || special)
}

func send_message(c *connection, h *hub, mesg *message_send){
	select {
	case c.send <- mesg:
	default:
		delete(h.connections, c)
		close(c.send)
	}
}

func (h *hub) run(name string ,clear chan string) {
	live := true
	for live {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			unregister(h, c)
			if will_self_destroy(len(h.connections), h.special){
				clear <-(name)
			}
		case m := <-h.broadcast:
			mesg := &message_send{Sender: m.id.name, Body:m.body}
			for c,_ := range h.connections {
				if c == m.id {
					continue
				} else {
					send_message(c, h, mesg)
			  }
			}
		case l := <-h.live:
			// close connection?
			live = l
		}
	}
}
