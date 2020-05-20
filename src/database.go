package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "os"
    "strconv"
    "strings"
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
    createDailySubscribersTableQuery = `CREATE TABLE IF NOT EXISTS usersDaily(
                                UserID SERIAL PRIMARY KEY,
                                ChatID INT,
                                userTimezone TEXT,
                                userMealsUTC TIME[]);`
    createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users(userID SERIAL PRIMARY KEY,
                                 chatID NUMERIC,
                                 username TEXT UNIQUE,
                                 patience INT DEFAULT 2,
                                 selectedFrequency INT DEFAULT 0,
                                 userTimezone INT DEFAULT -100,
                                 userMealEditIndex INT DEFAULT 0,
                                 userMealsUTC TIME[]);`
)

func selectIntValueFromUsers(valueType string, username string) (int, error) {
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return -100, err
    }
    defer db.Close()
    var query string
    query = "SELECT " + valueType + " FROM users WHERE username = '" + username + "';"
    log.Printf("SQL successfully initialized to execute query: " + query)
    row := db.QueryRow(query)
    value := 0
    err = row.Scan(&value)
    if err != nil {
        return -100, err
    }
    return value, nil
} 

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

func insertUser(username string) error {
    return execQuery("INSERT INTO users(username) VALUES ('" + username + "') ON CONFLICT (username) DO NOTHING;")
}

func clearInsertUser(username string) error {
    return execQuery("INSERT INTO users(username) VALUES ('" + username + "') ON CONFLICT (username) DO UPDATE SET patience = DEFAULT, selectedFrequency = DEFAULT, userTimezone = DEFAULT, userMealEditIndex = DEFAULT, userMealsUTC = DEFAULT;")
}

func addUserPatience(username string, patience int) error {
    return execQuery("INSERT INTO users(username, patience) VALUES ('" + username + "', " + strconv.Itoa(patience) + ") ON CONFLICT (username) DO UPDATE SET patience = users.patience + (" + strconv.Itoa(patience) + ");")
}

func setUserPatience(username string, patience int) error {
    return execQuery("INSERT INTO users(username, patience) VALUES ('" + username + "', " + strconv.Itoa(patience) + ") ON CONFLICT (username) DO UPDATE SET patience = " + strconv.Itoa(patience) + ";")
}

func createUserDailyMeals(username string) error {
    frequency, _ := getUserSelectedFrequency(username)
    emptyMeals := make([]string, frequency-1)
    for i := 1; i < frequency; i++ {
        emptyMeals[i-1] = "null"
    }
    meals := "[" + strings.Join(emptyMeals, ",") + "]"
    return execQuery("INSERT INTO users(username, userMealsUTC) VALUES ('" + username + "', ARRAY" + meals + "::TIME[]) ON CONFLICT (username) DO UPDATE SET userMealsUTC = ARRAY" + meals + "::TIME[];")
}

func updateUserDailyMeal(username string, mealTime string) error {
    a, _ := getUserMealEditIndex(username)
    log.Printf("userMealEditIndex =  : " + strconv.Itoa(a))
    return execQuery("UPDATE users SET userMealsUTC[" + "(SELECT userMealEditIndex FROM users WHERE username = '" + username + "')" + "-1] = '" + mealTime + "' WHERE username = '" + username + "';")
}

func userMealsUTCSet(username string) (bool, error) {
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return false, err
    }
    defer db.Close()
    var query string
    query = "SELECT array_position((SELECT userMealsUTC FROM users WHERE username = '" + username + "'), null);"
    log.Printf("SQL successfully initialized to execute query: " + query)
    row := db.QueryRow(query)
    value := -100
    err = row.Scan(&value)
    if value == -100 {
        return true, err
    } else {
        return false, err
    }
}

func userMealsUTCNotSet(username string) (*[]int, error) {
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return nil, err
    }
    defer db.Close()
    var query string
    query = "SELECT array_positions((SELECT userMealsUTC FROM users WHERE username = " + username + "), null);"
    log.Printf("SQL successfully initialized to execute query: " + query)
    row := db.QueryRow(query)
    var value []int
    err = row.Scan(&value)
    if err != nil {
        return nil, err
    }
    return &value, nil
}

func setUserSelectedFrequency(username string, frequency string) error {
    return execQuery("INSERT INTO users(username, selectedFrequency) VALUES ('" + username + "', " + frequency + ") ON CONFLICT (username) DO UPDATE SET selectedFrequency = " + frequency + ";")
}

func setUserTimezone(username string, timezone string) error {
    return execQuery("INSERT INTO users(username, userTimezone) VALUES ('" + username + "', " + timezone + ") ON CONFLICT (username) DO UPDATE SET userTimezone = " + timezone + ";")
}

func setUserMealEditIndex(username string, index string) error {
    return execQuery("INSERT INTO users(username, userMealEditIndex) VALUES ('" + username + "', " + index + ") ON CONFLICT (username) DO UPDATE SET userMealEditIndex = " + index + ";")
}


func getUserTimezone(username string) (int, error) {
    return selectIntValueFromUsers("userTimezone", username)
}

func getUserMealEditIndex(username string) (int, error) {
    return selectIntValueFromUsers("userMealEditIndex", username)
}

func getUserSelectedFrequency(username string) (int, error) {
     return selectIntValueFromUsers("selectedFrequency", username)
}

func getUserPatience(username string) (int, error) {
    return selectIntValueFromUsers("patience", username)
}

func createUsageTable() error {
    return execQuery(createUsageTableQuery)
}

func createDailyUsersTable() error {
    return execQuery(createDailySubscribersTableQuery)
}

func createUsersTable() error {
    return execQuery(createUsersTableQuery)
}