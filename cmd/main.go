package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

var db *sql.DB // GLOBAL DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	dataConnection()
	http.HandleFunc("/add", add)
	http.HandleFunc("/complete", complete)
	http.HandleFunc("/lists", list)
	http.HandleFunc("/delete", delete)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server running on :8080 ...")
	http.ListenAndServe(":8080", nil)
}

func dataConnection() {
	/***** DATABASE *****/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	query := `
		CREATE TABLE IF NOT EXISTS todolist (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
		    done BOOLEAN
		)
	`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table ready to use!")
	/*----- DATABASE -----*/
}

func add(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")

	_, err := db.Exec(`INSERT INTO todolist (name, done) VALUES ($1, false)`, name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Item %s inserted successfully!\n", name)
}

func complete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", 400)
		return
	}

	_, err = db.Exec(`UPDATE todolist SET done = true WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Completed id: %d", id)
}

func list(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, done FROM todolist ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var res []map[string]interface{}
	displayID := 1

	for rows.Next() {
		var dbID int
		var name string
		var done bool

		rows.Scan(&dbID, &name, &done)

		res = append(res, map[string]interface{}{
			"displayID": displayID, // Front uchun 1,2,3
			"id":        dbID,      // Asl database ID
			"name":      name,
			"done":      done,
		})

		displayID++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", 400)
		return
	}

	_, err = db.Exec(`DELETE FROM todolist WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Deleted id: %d", id)
}
