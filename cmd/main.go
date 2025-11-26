package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
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

	fmt.Printf("Item %s inserted successfully!\n", name)
}

func complete(w http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
	id_done, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", 400)
		return
	}

	update := (`
        UPDATE todolist
		SET done = true
		WHERE id = $1;
    `)

	_, err = db.Exec(update, id_done)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Item %d marked complete\n", id_done)
}

func list(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query(`SELECT id, name, done FROM todolist ORDER BY id ASC`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var done bool

		if err := rows.Scan(&id, &name, &done); err != nil {
			log.Println(err)
			continue
		}

		status := "-"
		if done {
			status = "+"
		}

		fmt.Fprintf(w, "%d) %s [%s]\n", id, name, status)
	}
}

func delete(w http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
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

	fmt.Fprintf(w, "Item %d deleted\n", id)
}
