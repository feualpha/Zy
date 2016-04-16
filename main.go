// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main

import (
  "crypto/md5"
  "database/sql"
  "encoding/hex"
  "flag"
  "github.com/goji/httpauth"
  "github.com/gorilla/mux"
  _"github.com/mattn/go-sqlite3"
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

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func myAuthFunc(username, password string) bool {
  db, err := sql.Open("sqlite3", "./foo.db")
  if err != nil {
    log.Fatal("error 201")
  }
  defer db.Close()

  query, err := db.Prepare("select password from foo where username = ?")
  if err != nil {
		log.Fatal("error 202")
	}
	defer query.Close()

  var q_password string
  err = query.QueryRow(username).Scan(&q_password)
	if err != nil {
    return false
	}

  hasher := md5.New()
  hasher.Write([]byte(password))
  hashed_pass := hex.EncodeToString(hasher.Sum(nil))  

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

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
