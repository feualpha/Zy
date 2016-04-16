package main

import (
  "encoding/json"
  "log"
  "net/http"
)

type registerRespond struct {
  status bool
  message string
}

type registerData struct {
  Username string
  Password string
}

func getRegisterValue(r *http.Request) (string, string){
  decoder := json.NewDecoder(r.Body)

	var t registerData
	err := decoder.Decode(&t)
	if err != nil {
    log.Println("gagal")
		log.Fatal(err)
	}

  return t.Username, encryptPass(t.Password)
}

func registerHandler(c http.ResponseWriter, r *http.Request){
  username, password :=  getRegisterValue(r)
  respond := registerRespond{status: true}
  err := dbRegister(username, password)
  if err != nil {
    respond.status = false
  }
  
  json.NewEncoder(c).Encode(respond)
}
