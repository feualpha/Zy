package main

import (
  "encoding/json"
  "log"
  "net/http"
)

const min_username_length int = 5
const min_password_length int = 5
const code_registration_success int =  400
const message_registration_success string =  "Registration Success"
const code_error_username_length int = 401
const message_error_username_length string =  "Error Minimum Username Length 5"
const code_error_password_length int = 402
const message_error_password_length string =  "Error Minimum Password Length 5"
const code_error_username_exist int = 403
const message_error_username_exist string =  "Error Username Already Exist"
const code_validate_success int = 404
const message_validate_success string =  "Validation Success"
const code_registration_failed int =  405
const message_registration_failed string =  "Registration Failed"


type registerRespond struct {
  Code int
  Message string
}

type registerData struct {
  Username string
  Password string
}

func validateRegistration(username, password string) int {
  b_username := []byte(username)
  b_password := []byte(password)
  if len(b_username) < min_username_length {
    return code_error_username_length
  } else if len(b_password) < min_password_length {
    return code_error_password_length
  } else if !checkUniqueness(username) {
    return code_error_username_exist
  } else {
    return code_validate_success
  }
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

func registering(username, password string) registerRespond {
  err := dbRegister(username, password)
  if err != nil {
    return respondComposer(code_registration_failed)
  }

  return respondComposer(code_registration_success)
}

func registerHandler(c http.ResponseWriter, r *http.Request){
  username, password :=  getRegisterValue(r)
  status := validateRegistration(username, password)
  var respond registerRespond

  if status == code_validate_success {
    respond = registering(username, password)
  } else {
    respond = respondComposer(status)
  }

  json.NewEncoder(c).Encode(respond)
}

func respondComposer(code int) registerRespond{
  var respond registerRespond
  switch code {
  case code_registration_success:
    respond = registerRespond{Code: code, Message: message_registration_success}
  case code_error_username_length:
    respond = registerRespond{Code: code, Message: message_error_username_length}
  case code_error_password_length:
    respond = registerRespond{Code: code, Message: message_error_password_length}
  case code_error_username_exist:
    respond = registerRespond{Code: code, Message: message_error_username_exist}
  case code_registration_failed:
    respond = registerRespond{Code: code, Message: message_registration_failed}
  }
  return respond
}
