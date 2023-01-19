package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []student `json:"data"`
	Message string    `json:"message"`
}

func setupDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	rout := mux.NewRouter()
	rout.HandleFunc("/", Get_stu).Methods("GET")
	rout.HandleFunc("/", Post_stu).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", rout))
}

func Get_stu(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	rows, err := db.Query("SELECT * FROM go_person")

	if err != nil {
		panic(err)
	}

	var stu []student

	for rows.Next() {
		var Id string
		var Nam string

		err = rows.Scan(&Id, &Nam)

		if err != nil {
			panic(err)
		}

		stu = append(stu, student{ID: Id, Name: Nam})
	}

	var response = JsonResponse{Type: "success", Data: stu}

	json.NewEncoder(w).Encode(response)
}

func Post_stu(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	var response = JsonResponse{}

	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO go_person(Name) VALUES(?)")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	json.Unmarshal(body, &keyVal)

	if keyVal["name"] == "" {
		response = JsonResponse{Type: "error", Message: "missing Name"}
	} else {
		_, err = stmt.Exec(keyVal["name"])

		if err != nil {
			panic(err.Error())
		}
		response = JsonResponse{Type: "success", Message: "Record inserted successfully!"}
	}

	response = JsonResponse{Type: "success", Message: "Record inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}
