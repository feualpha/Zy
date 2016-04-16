package main

import (
  "database/sql"
  _"github.com/mattn/go-sqlite3"
  "log"
)

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

  query, err := db.Prepare("select password from foo where username = ?")
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

  query, err := db.Prepare("INSERT INTO foo(username, password) values(?, ?)")
  if err != nil {
    log.Fatal("error 202")
  }
  defer query.Close()

  _, err = query.Exec(username, password)

  return err
}
