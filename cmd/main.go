package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	interactive := flag.Bool("i", false, "Interactive mode")
	flag.Parse()

	if *interactive {
		interactiveMode()
	}
}

func interactiveMode() {
	general := bufio.NewScanner(os.Stdin)

	/***** DATABASE *****/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

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

	for {
		fmt.Print("> ")

		if !general.Scan() {
			break
		}

		txt := general.Text()
		parts := strings.SplitN(txt, " ", 2)

		command := parts[0]

		switch command {
		case "add":
			var done bool

			fmt.Println("Item nomi: ")
			name := parts[1]

			qoshish := `INSERT INTO todolist (name, done) VALUES ($1, $2)`

			_, err = db.Exec(qoshish, name, done)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Item %s inserted successfully!\n", name)

		case "complete":
			fmt.Println("Item raqami: ")
			index := parts[1]
			id_done, errr := strconv.Atoi(index)
			if errr != nil {
				panic(errr)
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
			fmt.Println("Item updated to done:", id_done)

		case "lists":

			rows, err := db.Query(`select name, done from todolist ORDER BY id ASC`)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			counter := 1
			for rows.Next() {
				var name string
				var done bool
				var status string

				if err := rows.Scan(&name, &done); err != nil {
					log.Fatal(err)
				}

				if done == true {
					status = "+"
				} else if done == false {
					status = "-"
				}

				fmt.Printf("%d. [%s] %s\n", counter, status, name)
				counter++
			}

		case "delete":
			delete_name := parts[1]

			res, err := db.Exec(`delete from todolist where name=$1`, delete_name)
			if err != nil {
				log.Fatal(err)
			}

			affected, _ := res.RowsAffected()

			if affected == 0 {
				fmt.Println("Item topilmadi.")
			} else {
				fmt.Printf("Item %s deleted successfully!\n", delete_name)
			}

		default:
			fmt.Println("Write one of these first: -add, -complete, -lists, -delete.")
		}
	}
}
