package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dedeTr/url-shortener/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

func CreateURL(w http.ResponseWriter, r *http.Request) {

	// create an empty user of type models.User
	var linkDB models.LinkDB

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&linkDB)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertURL(linkDB)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func insertURL(linkDB models.LinkDB) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO linkdb (longurl, shorturl) VALUES ($1, $2) RETURNING urlid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, linkDB.LongURL, linkDB.ShortURL).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func Geturl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	url_req := params["url_req"]

	link, err := geturl(string(url_req))

	if err != nil {
		log.Fatalf("unable to get link, %v", err)
	}

	json.NewEncoder(w).Encode(link)
}

func geturl(short_url string) (models.LinkDB, error) {
	db := createConnection()

	defer db.Close()

	var linkdb models.LinkDB

	sqlQuery := `SELECT * FROM linkdb WHERE shorturl=$1`

	row := db.QueryRow(sqlQuery, short_url)

	err := row.Scan(&linkdb.Id, &linkdb.LongURL, &linkdb.ShortURL)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return linkdb, nil
	case nil:
		return linkdb, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return linkdb, err
}
