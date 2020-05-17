package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "os"
)

var host = os.Getenv("PGHOST")
var port = os.Getenv("PGPORT")
var user = os.Getenv("PGUSER")
var password = os.Getenv("PGPASSWORD")
var dbname = os.Getenv("PGDATABASE")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

const (
    createUsageTableQuery = `CREATE TABLE IF NOT EXISTS usage(
                                ID SERIAL PRIMARY KEY,
                                timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                username TEXT,
                                chatID INT,
                                message TEXT,
                                answer TEXT);`
    createDailyUsersTableQuery = `CREATE TABLE IF NOT EXISTS usersDaily(
                                UserID SERIAL PRIMARY KEY,
                                ChatID INT,
                                userTimezone TEXT,
                                userMealsUTC TIME[]);`
)


func execQuery(query string) error {
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return err
    }
    defer db.Close()
    log.Printf("SQL successfully initialized to execute query: " + query)

    _, err = db.Exec(query)

    if err != nil {
      return err
    }
    return nil
}

func createUsageTable() error {
    return execQuery(createUsageTableQuery)
}

func createDailyUsersTable() error {
    return execQuery(createDailyUsersTableQuery)
}