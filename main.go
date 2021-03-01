package main

import (
	"database/sql"
	_ "html/template"
	"log"
	"net/http"
	"text/template"

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

// regionHandler
type regionHandler struct {
	db         *sql.DB
	structItem []Item
}

// InitDB initializes SQLite Database
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./angebote.db")

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
		Laden TEXT);`)

	if err != nil {
		panic(err)
	}
}

// StoreItem inserts Items into database
func StoreItem(db *sql.DB, items []Item) {
	addItem, err := db.Prepare(`INSERT OR REPLACE INTO items(
		ID,
		Region,
		Kategorie,
		Angebot,
		Laden) values(?, ?, ?, ?, ?)`)

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
	ORDER BY ID ASC`)

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

// HomeHandler parses HTML files for the homepage
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	homeTemplate, _ := template.ParseFiles("static/index.html")
	err := homeTemplate.Execute(w, nil)

	if err != nil {
		log.Println("runtime error: exec template ", err)
		return
	}
}

// RegionHandler manages pages of selected items in requested Regions
func (h *regionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	regionString := req.URL.Query().Get(":reg")
	var regionResult []Item

	itemsForRegionHandler, _ := h.db.Query(`SELECT ID, Kategorie, Angebot, Laden FROM items
	WHERE Region LIKE '` + regionString + `' ORDER BY ID ASC`)

	defer itemsForRegionHandler.Close()
	for itemsForRegionHandler.Next() {
		item := Item{}
		itemsForRegionHandler.Scan(&item.ID, &item.Kategorie, &item.Angebot, &item.Laden)
		regionResult = append(regionResult, item)
	}

	regionData := map[string]interface{}{
		"Region":      regionString,
		"RegionItems": regionResult,
	}

	if regionString == "Nord" || regionString == "West" || regionString == "Süd" {
		regionTemplate, _ := template.ParseFiles("static/region.html")
		regionTemplate.Execute(w, regionData)
	} else {
		badTemplate, _ := template.ParseFiles("static/404.html")
		badTemplate.Execute(w, regionData)
	}

}

// ListHandler parses file for page of saved items
func ListHandler(w http.ResponseWriter, req *http.Request) {
	listTemplate, _ := template.ParseFiles("static/merkliste.html")
	listTemplate.Execute(w, nil)
}

// FormHandler listens for Form to post new items
func FormHandler(w http.ResponseWriter, req *http.Request) {
	listTemplate, _ := template.ParseFiles("static/formular.html")
	// listener
	// StoreItem(db, formItem)
	listTemplate.Execute(w, nil)
}

// gets Values from Item Form 
func form(w http.ResponseWriter, r *http.Request){

    if r.Method != "GET" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
		return         
    }

    PId := r.FormValue("id")
	PReg := r.FormValue("Region")
	PKat := r.FormValue("Kategorie")
	PAng := r.FormValue("Angebot")
	PLad := r.FormValue("Laden")

    
	formTemplate, _ := template.ParseFiles("static/formular.html")
	formTemplate.Execute(w, nil)

	//tpl.ExecuteTemplate(w, "formular.html")
}

func main() {
	db := InitDB()
	defer db.Close()
	CreateTable(db)

	testItems := []Item{
		{"1", "Nord", "Essen", "Pizzas dienstags für 7 Euro", "Pizzeria XY in Hünfeld"},
		{"2", "West", "Kleidung", "Sale bis 50% auf T-Shirts", "Klamottenladen in Fulda"},
		{"3", "Nord", "Technik", "USB-C Kabel für nur 2,99€", "Tech Shop in Petersberg"},
		{"4", "West", "Essen", "Große Waffeln - 4€", "Waffelladen in Fulda"},
		{"5", "West", "Kleidung", "10€ Rabatt auf alle Jacken", "Klamottenladen 2 in Fulda"},
	}

	StoreItem(db, testItems)
	readItems := ReadItem(db)
	log.Println(readItems)

	m := pat.New()
	m.Get("/", http.HandlerFunc(HomeHandler))
	m.Get("/region/:reg", &regionHandler{db, testItems})

	/* fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs) */

	http.Handle("/", m)
	http.Handle("/merkliste", http.HandlerFunc(ListHandler))
	http.Handle("/formular", http.HandlerFunc(FormHandler))
	http.HandleFunc("/form",form)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
