// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type mesg struct {
	body []byte
	id *connection
}

type client struct {
	id *connection
	race bool
}

type hub struct {
	connections map[*connection]bool
	broadcast   chan *mesg
	register    chan *client
	unregister  chan *connection
	live        chan bool
}

func newHub() *hub {
	return &hub{
		broadcast:   make(chan *mesg),
		register:    make(chan *client),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
		live: make(chan bool),
	}
}

func (h *hub) run() {
	live := true
	for live {
		select {
		case c := <-h.register:
			h.connections[c.id] = c.race
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c, v := range h.connections {
				if v && (c == m.id) {
					continue
				}
				select {
				case c.send <- m.body:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		case l := <-h.live:
			// close connection?
			live = l
		}
	}
}
