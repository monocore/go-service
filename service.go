package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Movie Struct
type Movie struct {
	Title  string `json:"title"`
	Rating string `json:"rating"`
	Year   string `json:"year"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/movies", handleMovies).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func getJSON(connection string, sqlString string) (string, error) {
	//"root@/movies"
	db, error := sql.Open("mysql", connection)
	if error != nil {
		log.Println(error.Error())
		return "", error
	}

	error = db.Ping()
	if error != nil {
		log.Println(error.Error())
		return "", error
	}
	// Prepare statement for reading data
	stmtOut, error := db.Prepare(sqlString)
	if error != nil {
		log.Println(error.Error())
		return "", error
	}
	defer stmtOut.Close()

	// Execute the query
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Println(error.Error())
		return "", error
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Println(error.Error())
		return "", error
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

func handleMovies(res http.ResponseWriter, req *http.Request) {
	var json, error = getJSON("root@/movies", "select * from movies")
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(res, json)
	return
}
