package main
//change run hub so that it can kill it self
//switcher routine return create room, return desired room,
//destory room if empty default room if not specified

type switch_agent struct{
  id *string
  room chan *hub
}

type switcher struct {
  rooms map[string] *hub
  join chan switch_agent
}

func newSwitcher() *switcher {
	return &switcher{
		rooms: make(map[string]*hub),
    join: make(chan switch_agent),
	}
}

func (s *switcher) run() {
  s.rooms["default"] = newHub()
  for {
    select {
    case c:= <-s.join:
      hub, ok := s.rooms[*c.id]
      if ok {
        c.room <- hub
      } else {
        c.room <-s.rooms["default"]
      }
    }
  }
}
