package main

import (
  "database/sql"
  _"github.com/mattn/go-sqlite3"
  "log"
)

const query_select string = "select password from foo where username = ?"
const query_insert string = "INSERT INTO foo(username, password) values(?, ?)"

func openDb() *sql.DB {
  db, err := sql.Open("sqlite3", "./foo.db")
  if err != nil {
    log.Fatal("error 201")
  }

  return db
}

func dbAuth(username string) string {
  q_password := ""

  db := openDb()
  defer db.Close()

  query, err := db.Prepare(query_select)
  if err != nil {
    log.Fatal("error 202")
  }
  defer query.Close()

  _ = query.QueryRow(username).Scan(&q_password)

	return q_password
}

func dbRegister(username, password string) error {
  db := openDb()
  defer db.Close()

  query, err := db.Prepare(query_insert)
  if err != nil {
    log.Fatal("error 202")
  }
  defer query.Close()

  _, err = query.Exec(username, password)

  return err
}

func checkUniqueness(username string) bool {
  db := openDb()
  defer db.Close()

  query, err := db.Prepare(query_select)
  if err != nil {
    log.Fatal("error 202")
  }
  defer query.Close()

  temp := ""
  err = query.QueryRow(username).Scan(&temp)
  if err != nil{
    return true
  }

  return false
}
