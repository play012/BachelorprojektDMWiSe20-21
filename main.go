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
	_, err := db.Exec(`DROP TABLE IF EXISTS items; CREATE TABLE items(
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
	req.ParseForm()
	sortReq := req.Form.Get("sort")

	//newest items and page without sort request serve the same items
	if len(sortReq) == 0 || sortReq == "newest" {
		newestItems, _ := h.db.Query(`SELECT ID, Kategorie, Angebot, Laden FROM items
		WHERE Region LIKE '` + regionString + `' ORDER BY ID DESC`)

		defer newestItems.Close()
		for newestItems.Next() {
			item := ReadItem{}
			newestItems.Scan(&item.ID, &item.Kategorie, &item.Angebot, &item.Laden)
			regionResult = append(regionResult, item)
		}

		sortReq = "newest"
		regionData := map[string]interface{}{
			"Region":      regionString,
			"RegionItems": regionResult,
			"Sort":        sortReq,
		}

		if regionString == "Nord" || regionString == "West" || regionString == "Süd" {
			regionTemplate, _ := template.ParseFiles("static/region.html")
			regionTemplate.Execute(w, regionData)
		} else {
			badTemplate, _ := template.ParseFiles("static/404.html")
			badTemplate.Execute(w, nil)
		}
	} else if sortReq == "oldest" {
		// oldest items serves items with lowest ID first
		oldestItems, _ := h.db.Query(`SELECT ID, Kategorie, Angebot, Laden FROM items
		WHERE Region LIKE '` + regionString + `' ORDER BY ID ASC`)

		defer oldestItems.Close()
		for oldestItems.Next() {
			item := ReadItem{}
			oldestItems.Scan(&item.ID, &item.Kategorie, &item.Angebot, &item.Laden)
			regionResult = append(regionResult, item)
		}

		regionData := map[string]interface{}{
			"Region":      regionString,
			"RegionItems": regionResult,
			"Sort":        sortReq,
		}

		if regionString == "Nord" || regionString == "West" || regionString == "Süd" {
			regionTemplate, _ := template.ParseFiles("static/region.html")
			regionTemplate.Execute(w, regionData)
		} else {
			badTemplate, _ := template.ParseFiles("static/404.html")
			badTemplate.Execute(w, regionData)
		}
	} else {
		badTemplate, _ := template.ParseFiles("static/404.html")
		badTemplate.Execute(w, nil)
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

		SaveItem(h.db, []StoreItem{{pReg, pKat, pAng, pLad}})
		// Test
		addItem, err := h.db.Prepare(`INSERT OR REPLACE INTO items(
			Region,
			Kategorie,
			Angebot,
			Laden) VALUES (?, ?, ?, ?);`)

		if err != nil {
			panic(err)
		}
		defer addItem.Close()

		//addItem.Exec(&pReg, &pKat, &pAng, &pLad)
		//addItem.Exec("Süd", "Essen", "Test", "Test")

		rows, _ := h.db.Query("SELECT Region, Kategorie, Angebot, Laden FROM items")
		var region string
		var kategorie string
		var angebot string
		var laden string
		for rows.Next() {
			rows.Scan(&region, &kategorie, &angebot, &laden)
			//log.Println(region + "," + kategorie + "," + angebot + "," + laden)
		}

		//

		log.Println("Neues Item gespeichert: "+pReg, pKat, pAng, pLad)
		formTemplate.Execute(w, struct{ Success bool }{true})
	} else {
		formTemplate.Execute(w, nil)
	}
}

// ImpressumHandler serves impressum.html
func ImpressumHandler(w http.ResponseWriter, req *http.Request) {
	impressumTpl, _ := template.ParseFiles("static/impressum.html")
	impressumTpl.Execute(w, nil)
}

// NotFoundHandler catches requests for nonexistent routes and redirects to a 404 page
func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	notFoundTemplate, _ := template.ParseFiles("static/404.html")
	notFoundTemplate.Execute(w, nil)
}

//var db *sql.DB

func main() {
	db := InitDB()
	defer db.Close()
	CreateTable(db)

	testItems := []StoreItem{
		{"Nord", "Essen", "Pizzas dienstags für 7 Euro", "Pizzeria XY in Hünfeld"},
		{"West", "Kleidung", "Sale bis 50% auf T-Shirts", "Klamottenladen in Fulda"},
		{"Nord", "Technik", "USB-C Kabel für nur 2,99€", "Tech Shop in Petersberg"},
		{"West", "Essen", "Große Waffeln - 4€", "Waffelladen in Fulda"},
		{"Süd", "Kleidung", "50€ Rabatt auf alle Jacken", "Klamottenladen 2 in Eichenzell"},

		{"Nord", "Kleidung", "Kinder-Jeans für 15€", "Klamottenladen in Eiterfeld"},
		{"Nord", "Kleidung", "Bio-Pullover für 30€", "Klamottenladen in Burghaun"},
		{"Nord", "Technik", "Playstation V für 500€", "Game Laden in Nüsttal"},
		{"Nord", "Technik", "15% Rabatt auf Headphones", "Tech Shop in Hilders"},
		{"Nord", "Essen", "Selbstgemachtes Hof Eis 0,90€", "Hof-Laden in Eiterfeld"},

		{"Süd", "Kleidung", "T-Shirts 2-Pack 12,99€", "Klammottenladen in Dipperz"},
		{"Süd", "Technik", "Notebooks 10% Rabatt", "Tech Shop Künzell"},
		{"Süd", "Essen", "Pommes diesen Monat nur 2€", "Imbiss Shop Fulda"},
		{"Süd", "Kleidung", "15% Rabatt auf Graue Jeans", "Laden in Fulda"},
		{"Süd", "Essen", "Burger ab 5€", "Laden in Fulda"},
		{"Süd", "Kleidung", "Marken-Shirt 18€", "Laden in Fulda"},

		{"West", "Technik", "Bose Soundbox 199€", "Tech Shop in Hosenfeld"},
		{"West", "Kleidung", "Socken 5er-Pack 10€", "Klammottenladen in Neuhof"},
		{"West", "Essen", "Sandwichs 2,50€", "Imbiss Shop in Flieden"},
		{"West", "Technik", "Keyboard 17,50€", "Laden in Kalbach"},
	}

	SaveItem(db, testItems)
	/*
		//TEST
		items := []StoreItem{
			{"Süd", "Essen", "1", "2"},
		}
		addItem, erro := db.Prepare(`INSERT OR REPLACE INTO items(
			Region,
			Kategorie,
			Angebot,
			Laden) VALUES (?, ?, ?, ?)`)

		if erro != nil {
			panic(erro)
		}
		defer addItem.Close()
		for _, item := range items {
			addItem.Exec(item.Region, item.Kategorie, item.Angebot, item.Laden)
		}

		// ----------------

	*/
	m := pat.New()
	m.NotFound = http.HandlerFunc(NotFoundHandler)
	m.Get("/", http.HandlerFunc(HomeHandler))
	m.Get("/region/:reg", &RegionHandler{db, ShowItem(db)})
	m.Get("/merkliste", http.HandlerFunc(ListHandler))
	m.Get("/formular", &FormHandler{db})
	m.Post("/formular", &FormHandler{db})
	m.Get("/impressum", http.HandlerFunc(ImpressumHandler))

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
