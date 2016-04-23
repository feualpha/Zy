package main
//change run hub so that it can kill it self
//switcher routine return create room, return desired room,
//destory room if empty default room if not specified
//import "log"

type switch_agent struct{
  id string
  room chan *hub
}

type switcher struct {
  rooms map[string] *hub
  join chan switch_agent
  create chan switch_agent
  clear chan string
}

func newSwitchAgent(id string) *switch_agent{
  return &switch_agent{
		id: id,
    room: make(chan *hub),
	}
}

func newSwitcher() *switcher {
	return &switcher{
		rooms: make(map[string]*hub),
    join: make(chan switch_agent),
    create: make(chan switch_agent),
    clear: make(chan string),
	}
}

func (s *switcher) run() {
  d_name := "default"
  first := newHub(true)
  s.rooms[d_name] = first
  go first.run(d_name, s.clear)
  for {
    select {
    case c := <-s.join:
      hub, ok := s.rooms[c.id]
      if ok {
        c.room <- hub
      } else {
        c.room <-s.rooms["default"]
      }
    case c := <-s.create:
      new_hub := newHub(false)
      s.rooms[c.id] = new_hub
      go new_hub.run(d_name, s.clear)
      c.room <- new_hub
    case key := <-s.clear:
      hub := s.rooms[key]
      hub.live <- false
      delete(s.rooms, key)
    }
  }
}
