package db

import (
	"database/sql"
	"fmt"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// Page stores data for subpages
type Page struct {
	ID    int64
	Title string
	Body  []byte
}

// Frontpage stores slice of strings of pages titles
type Frontpage struct {
	Titles []string
}

// Database contains mysql database variable
var Database *sql.DB

// Initialize connect to the database
func Initialize() error {
	var err error
	Database, err = sql.Open("mysql", "root:roma@tcp(127.0.0.1:3306)/wiki") // TO DO PASSWORD TO DATABASE AS ARGUMENT
	if err != nil {
		return fmt.Errorf("error during opening database: %v", err)
	}

	err = Database.Ping()
	if err != nil {
		return fmt.Errorf("error during verifying connection to the database: %v", err)
	}
	fmt.Println("DB connected!")

	return nil
}

// Close closes the database
func Close() {
	Database.Close()
}

// InsertPage make insert query to store data of page in db
func InsertPage(p *Page) (int64, error) {
	result, err := Database.Exec("INSERT INTO pages (title, body) VALUES (?, ?)", p.Title, p.Body)
	if err != nil {
		return -1, fmt.Errorf("saving page %s to db: %v", p.Title, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("getting id generated by the database: %v", err)
	}
	return id, err
}

// UpdatePage make update query to edit existing page in database
func UpdatePage(p *Page) error {
	_, err := Database.Exec("UPDATE pages SET body = ? WHERE title = ?", p.Body, p.Title)
	if err != nil {
		return fmt.Errorf("updating page %s to db: %v", p.Title, err)
	}
	return nil
}

// GetTitles make query for titles of last created 10 pages
func GetTitles() (*Frontpage, error) {
	rows, err := Database.Query("SELECT title FROM pages ORDER BY pageId DESC LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("")
	}
	defer rows.Close()
	slice := make([]string, 0, 10)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("error during scan each row: %v", err)
		}
		slice = append(slice, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error from the overall query: %v", err)
	}

	return &Frontpage{Titles: slice}, nil
}

// PageByTitle make query for a single row from pages
func PageByTitle(title string) (*Page, error) {
	var p Page
	row := Database.QueryRow("SELECT * FROM pages WHERE title = ?", title)
	if err := row.Scan(&p.ID, &p.Title, &p.Body); err != nil {
		return &p, fmt.Errorf("page title:%s scan error: %v", title, err)
	}
	return &p, nil
}
