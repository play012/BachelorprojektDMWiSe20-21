package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	_ "github.com/mattn/go-sqlite3"
)

// Item is listed in Database
type Item struct {
	ID        string
	Region    string
	Kategorie string
	Angebot   string
	Laden     string
}

// InitDB initializes SQLite Database
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}

	return db
}

// CreateTable if not exists
func CreateTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS items(
		ID TEXT NOT NULL PRIMARY KEY,
		Region TEXT,
		Kategorie TEXT,
		Angebot TEXT,
		Laden TEXT,
		InsertedDatetime DATETIME);`)

	if err != nil {
		panic(err)
	}
}

// StroItem inserts Items into database
func StoreItem(db *sql.DB, items []Item) {
	addItem, err := db.Prepare(`INSERT OR REPLACE INTO items(
		ID,
		Region,
		Kategorie,
		Angebot,
		Laden,
		InsertedDatetime
	) values(?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`)

	if err != nil {
		panic(err)
	}
	defer addItem.Close()

	for _, item := range items {
		_, err2 := addItem.Exec(item.ID, item.Region, item.Kategorie, item.Angebot, item.Laden)
		if err2 != nil {
			panic(err2)
		}
	}
}

// ReadItem selects an item from database
func ReadItem(db *sql.DB) []Item {
	rows, err := db.Query(`SELECT ID, Region, Kategorie, Angebot, Laden FROM items
	ORDER BY datetime(InsertedDatetime) DESC`)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Item
	for rows.Next() {
		item := Item{}
		err2 := rows.Scan(&item.ID, &item.Region, &item.Kategorie, &item.Angebot, &item.Laden)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

// RegionHandler manages pages of selected items in requested Regions
func RegionHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, req, "static/region.html")
}

// ListHandler shows saved Items in a list
func ListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, req, "static/merkliste.html")
}

func main() {
	m := pat.New()
	m.Get("/:region", http.HandlerFunc(RegionHandler))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.Handle("/:region", m)
	http.Handle("/merkliste", http.HandlerFunc(ListHandler))

	db := InitDB()
	defer db.Close()
	CreateTable(db)

	testItems := []Item{
		Item{"1", "Nord", "Essen", "Pizzas dienstags für 7 Euro", "Pizzeria XY in Hünfeld"},
		Item{"2", "West", "Kleidung", "Sale bis 50% auf T-Shirts", "Klamottenladen in Fulda"},
	}

	StoreItem(db, testItems)
	readItems := ReadItem(db)
	log.Println(readItems)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
