package main

import (
  "database/sql"
  "log"
  _"github.com/mattn/go-sqlite3"
  "os"
)

const query_select string = "select password from users where username = ?"
const query_insert string = "INSERT INTO users(username, password) values(?, ?)"

func openDb() *sql.DB {
  db, err := sql.Open("sqlite3", os.Getenv("DB_NAME"))
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

func dbInit(){
  db := openDb()
  defer db.Close()
  sql_table := `
  CREATE TABLE IF NOT EXISTS users(
		id INTEGER NOT NULL PRIMARY KEY,
		username STRING,
		password STRING
	);
  `
  _, err := db.Exec(sql_table)
	if err != nil { log.Fatal("error 203") }
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
