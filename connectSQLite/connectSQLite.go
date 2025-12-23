package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbname = "ch06.db"

func insertData(db *sql.DB, dsc string) error {
	cT := time.Now().Format(time.RFC1123)
	stmt, err := db.Prepare("INSERT INTO book VALUES(NULL,?,?);")
	if err != nil {
		fmt.Println("Insert data table: ", err)
		return err
	}

	_, err = stmt.Exec(cT, dsc)
	if err != nil {
		fmt.Println("Insert data table:", err)
		return err
	}
	return nil
}

func selectData(db *sql.DB, n int) error {
	rows, err := db.Query("Select * from book WHERE id > ? ", n)
	if err != nil {
		fmt.Println("Select:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var dt string
		var description string

		err = rows.Scan(&id, &dt, &description)
		if err != nil {
			fmt.Println("Row: ", err)
			return err
		}
	}

	date, err := time.Parse(time.RFC1123, dt)
	if err != nil {
		fmt.Println("Date:", err)
		return err
	}
	fmt.Printf("%d %s %s\n", id, date, description)
	return nil
}

func main() {
	os.Remove(dbname)

	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer db.Close()

	const create string = `
	CREATE TABLE IF NOT EXIST book (
		id INTEGER NOT NULL PRIMARY KEY,
		time TEXT NOT NULL,
		description TEXT);`

	_, err = db.Exec(create)
	if err != nil {
		fmt.Println("Create table:", err)
		return
	}

	for i := 1; i < 11; i = i + 1 {
		dcs := "Description" + strconv.Itoa(i)
		err = insertData(db, dcs)
		if err != nil {
			fmt.Println("Insert data:", err)
		}
	}

	err = selectData(db, 5)
	if err != nil {
		fmt.Println("Select:", err)
	}

	time.Sleep(time.Second)
	cT := time.Now().Format(time.RFC1123)
	db.Exec("UPDATE book SET time = ? WHERE id > ?", cT, 7)

	err = selectData(db, 8)
	if err != nil {
		fmt.Println("Select:", err)
		return
	}
	stmt, err := db.Prepare("DELETE from book where id = ?")
	_, err = stmt.Exec(8)
	if err != nil {
		fmt.Println("Delete:", err)
		return
	}

	err = selectData(db, 7)
	if err != nil {
		fmt.Println("Select:", err)
		return
	}

	query, err := db.Query("SELECT count(*) as count from book")
	if err != nil {
		fmt.Println("Select:", err)
		return
	}
	defer query.Close()
	count := -100
	for query.Next() {
		_ = query.Scan(&count)
	}
	fmt.Println("count(*):", count)
}
