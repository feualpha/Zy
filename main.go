// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
  "crypto/md5"
  "encoding/hex"
  "flag"
  "github.com/goji/httpauth"
  "github.com/gorilla/mux"
  "go/build"
  "log"
  "net/http"
	"path/filepath"
	"text/template"
)

var (
	addr      = flag.String("addr", ":8080", "server address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func encryptPass(password string) string {
  hasher := md5.New()
  hasher.Write([]byte(password))
  return hex.EncodeToString(hasher.Sum(nil))
}

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func myAuthFunc(username, password string) bool {
  q_password := dbAuth(username)
  hashed_pass := encryptPass(password)

  return q_password == hashed_pass
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	h := newHub()
	go h.run()

  authOpts := httpauth.AuthOptions{ AuthFunc: myAuthFunc }

  r := mux.NewRouter()
  r.HandleFunc("/", homeHandler)
  r.Handle("/ws", wsHandler{h: h, race:false})
  r.Handle("/wsc", wsHandler{h: h, race:true})
  http.Handle("/", httpauth.BasicAuth(authOpts)(r))
  http.HandleFunc("/cregister", registerHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
