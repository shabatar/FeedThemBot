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
                                chatID BIGINT,
                                message TEXT,
                                answer TEXT);`
	createDailySubscribersTableQuery = `CREATE TABLE IF NOT EXISTS usersDaily(
                                userID SERIAL PRIMARY KEY,
                                chatID BIGINT,
                                username TEXT UNIQUE,
                                userMealsUTC TEXT[]);`
	createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users(userID SERIAL PRIMARY KEY,
                                 chatID BIGINT,
                                 username TEXT UNIQUE,
                                 patience INT DEFAULT 2,
                                 selectedFrequency INT DEFAULT 0,
                                 userTimezone INT DEFAULT -100,
                                 userMealEditIndex INT DEFAULT -100,
                                 userMealsUTC TEXT[]);`
)

type User struct {
	chatID            int64
	username          string
	patience          int
	selectedFrequency int
	userTimezone      int
	userMealEditIndex int
	userMealsUTC      []string
}

func getUserData(username string) (*User, error) {
	user := new(User)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return user, err
	}
	defer db.Close()

	user.selectedFrequency, err = getUserSelectedFrequency(username)
	user.userMealEditIndex, err = getUserMealEditIndex(username)
	user.patience, err = getUserPatience(username)
	user.userTimezone, err = getUserTimezone(username)
	user.userMealsUTC, err = selectStringArrayValueFromUsers("SELECT UNNEST(userMealsUTC) FROM users WHERE username = '" + username + "';")
	user.chatID, _ = getUserChatID(username)

	if err != nil {
		log.Printf("an error occured")
	}
	return user, nil
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

func insertUser(username string, chatID int64) error {
	return execQuery("INSERT INTO users(username, chatID) VALUES ('" + username + "', " + strconv.FormatInt(chatID, 10) + "::BIGINT) ON CONFLICT (username) DO UPDATE SET chatID = " + strconv.FormatInt(chatID, 10) + "::BIGINT;")
}

func migrateDailyUser(username string) error {
	return execQuery("INSERT INTO usersDaily(userid, chatid, username, usermealsutc) SELECT userid, chatid, username, usermealsutc FROM users WHERE username = '" + username + "' ON CONFLICT(username) DO UPDATE SET userid = excluded.userid, chatid = excluded.chatid, usermealsutc = excluded.usermealsutc;")
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
	for i := 0; i < frequency-1; i++ {
		emptyMeals[i] = "null"
	}
	meals := "[" + strings.Join(emptyMeals, ",") + "]"
	return execQuery("INSERT INTO users(username, userMealsUTC) VALUES ('" + username + "', ARRAY" + meals + "::TEXT[]) ON CONFLICT (username) DO UPDATE SET userMealsUTC = ARRAY" + meals + "::TEXT[];")
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

func setUserSelectedFrequency(username string, frequency string) error {
	return execQuery("INSERT INTO users(username, selectedFrequency) VALUES ('" + username + "', " + frequency + ") ON CONFLICT (username) DO UPDATE SET selectedFrequency = " + frequency + ";")
}

func setUserTimezone(username string, timezone string) error {
	return execQuery("INSERT INTO users(username, userTimezone) VALUES ('" + username + "', " + timezone + ") ON CONFLICT (username) DO UPDATE SET userTimezone = " + timezone + ";")
}

func setUserMealEditIndex(username string, index string) error {
	return execQuery("INSERT INTO users(username, userMealEditIndex) VALUES ('" + username + "', " + index + ") ON CONFLICT (username) DO UPDATE SET userMealEditIndex = " + index + ";")
}

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

func selectInt64ValueFromUsers(valueType string, username string) (int64, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return -100, err
	}
	defer db.Close()
	var query string
	query = "SELECT " + valueType + " FROM users WHERE username = '" + username + "';"
	log.Printf("SQL successfully initialized to execute query: " + query)
	row := db.QueryRow(query)
	var value int64
	err = row.Scan(&value)
	if err != nil {
		return -100, err
	}
	return value, nil
}

func selectStringArrayValueFromUsers(query string) ([]string, error) {
	var result []string
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return result, err
	}
	defer db.Close()
	log.Printf("SQL successfully initialized to execute query: " + query)
	rows, err := db.Query(query)
	for rows.Next() {
		var elem string
		err := rows.Scan(&elem)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, elem)
	}
	namesstr := "[" + strings.Join(result, ",") + "]"
	log.Printf("username " + namesstr)
	if err != nil {
		return result, err
	}
	return result, nil
}

func getClosestDailyUsers() ([]string, error) {
	var names []string
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return names, err
	}
	defer db.Close()
	var query string
	query = `with joinedSchedule as 
                (with schedule as 
                    (select username, unnest(usermealsutc) as usertime from usersDaily order by usertime)
                    select username, (DATE_PART('hour', usertime::time - now()::time))*60 + DATE_PART('minute', usertime::time - now()::time) as diff from schedule)
            select username from joinedSchedule where diff > 0 and diff < 16;`
	log.Printf("SQL successfully initialized to execute query: " + query)
	rows, err := db.Query(query)
	name := ""
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}
	namesstr := "[" + strings.Join(names, ",") + "]"
	log.Printf("username " + namesstr)
	if err != nil {
		return names, err
	}
	return names, nil
}

func syncTimezone(username string) error {
	query := `update users set usermealsutc = (select array(select mealtime::time - '01:00:00'::time * usertimezone from (select usertimezone, unnest(usermealsutc) as mealtime from users where username = '` + username + `') as t1)::text[]) where username = '` + username + `';`
	return execQuery(query)
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

func getUserChatID(username string) (int64, error) {
	return selectInt64ValueFromUsers("chatID", username)
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
