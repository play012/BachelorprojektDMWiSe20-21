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

// ReadItem is the item listed in database with its ID
type ReadItem struct {
	ID        int
	Region    string
	Kategorie string
	Angebot   string
	Laden     string
}

// StoreItem is the item to store in the database (without ID -> autoincrement)
type StoreItem struct {
	Region    string
	Kategorie string
	Angebot   string
	Laden     string
}

// RegionHandler gives database and Items read to regionpage
type RegionHandler struct {
	db         *sql.DB
	structItem []ReadItem
}

// FormHandler gives database to form for adding new items
type FormHandler struct {
	db *sql.DB
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
	// AUTOINCREMENT creates own IDs as primary keys
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS items(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Region TEXT,
		Kategorie TEXT,
		Angebot TEXT,
		Laden TEXT);`)

	if err != nil {
		panic(err)
	}
}

// SaveItem inserts Items into database
func SaveItem(db *sql.DB, items []StoreItem) {
	addItem, err := db.Prepare(`INSERT OR REPLACE INTO items(
		Region,
		Kategorie,
		Angebot,
		Laden) values(?, ?, ?, ?)`)

	if err != nil {
		panic(err)
	}
	defer addItem.Close()

	for _, item := range items {
		addItem.Exec(item.Region, item.Kategorie, item.Angebot, item.Laden)
	}
}

// ShowItem selects an item from database
func ShowItem(db *sql.DB) []ReadItem {
	rows, err := db.Query(`SELECT ID, Region, Kategorie, Angebot, Laden FROM items
	ORDER BY ID ASC`)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []ReadItem
	for rows.Next() {
		item := ReadItem{}
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
		log.Println("runtime error: exec home template ", err)
		return
	}
}

// RegionHandler manages pages of selected items in requested Regions
func (h *RegionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	regionString := req.URL.Query().Get(":reg")
	var regionResult []ReadItem

	itemsForRegionHandler, _ := h.db.Query(`SELECT ID, Kategorie, Angebot, Laden FROM items
	WHERE Region LIKE '` + regionString + `' ORDER BY ID ASC`)

	defer itemsForRegionHandler.Close()
	for itemsForRegionHandler.Next() {
		item := ReadItem{}
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



// FormHandler gets values from Item Form
func (h *FormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formTemplate, _ := template.ParseFiles("static/formular.html")

	if r.Method == http.MethodPost {
		r.ParseForm()
		pReg := r.FormValue("Region")
		pKat := r.FormValue("Kategorie")
		pAng := r.FormValue("Angebot")
		pLad := r.FormValue("Laden")

		newItems := []StoreItem{
			{"Süd", "Essen", "Pizzas dienstags für 7 Euro", "Pizzeria XY in Hünfeld"},
			{pReg, pKat, pAng, pLad},
		}

		SaveItem(h.db, newItems)
		
		log.Println("Neues Item gespeichert: "+pReg, pKat, pAng, pLad)
		formTemplate.Execute(w, struct{ Success bool }{true})
	} else {
		formTemplate.Execute(w, nil)
	}
}

// NotFoundHandler catches requests for nonexistent routes and redirects to a 404 page
func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	notFoundTemplate, _ := template.ParseFiles("static/404.html")
	notFoundTemplate.Execute(w, nil)
}

func main() {
	db := InitDB()
	defer db.Close()
	CreateTable(db)

	testItems := []StoreItem{
		{"Nord", "Essen", "Pizzas dienstags für 7 Euro", "Pizzeria XY in Hünfeld"},
		{"West", "Kleidung", "Sale bis 50% auf T-Shirts", "Klamottenladen in Fulda"},
		{"Nord", "Technik", "USB-C Kabel für nur 2,99€", "Tech Shop in Petersberg"},
		{"West", "Essen", "Große Waffeln - 4€", "Waffelladen in Fulda"},
		{"West", "Kleidung", "10€ Rabatt auf alle Jacken", "Klamottenladen 2 in Fulda"},
	}

	SaveItem(db, testItems)

	m := pat.New()
	m.NotFound = http.HandlerFunc(NotFoundHandler)
	m.Get("/", http.HandlerFunc(HomeHandler))
	m.Get("/region/:reg", &RegionHandler{db, ShowItem(db)})
	m.Get("/merkliste", http.HandlerFunc(ListHandler))
	m.Get("/formular", &FormHandler{db})
	m.Post("/formular", &FormHandler{db})

	/* fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs) */

	http.Handle("/", m)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("RUNNING")

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
