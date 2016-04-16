// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main

import (
  "net/http"
	"flag"
	"go/build"
	"log"
	"path/filepath"
	"text/template"
  "github.com/gorilla/mux"
  "github.com/abbot/go-http-auth"
)

var (
	addr      = flag.String("addr", ":8080", "server address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func Secret(user, realm string) string {
  if user == "john" {
    // password is "hello"
    return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
    }
  return ""
}

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	h := newHub()
	go h.run()
  /////////////
  r := mux.NewRouter()

  r.HandleFunc("/", homeHandler)
  r.Handle("/ws", wsHandler{h: h, race:false})
  r.Handle("/wsc", wsHandler{h: h, race:true})
  http.Handle("/", httpauth.SimpleBasicAuth("dave", "somepassword")(r))
  ////////////////
	//http.HandleFunc("/", homeHandler)
	//http.Handle("/ws", wsHandler{h: h, race:false})
  //http.Handle("/wsc", wsHandler{h: h, race:true})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
