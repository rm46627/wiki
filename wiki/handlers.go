package wiki

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/rm46627/wiki/db"
)

// Handler handles requests for front and notfound pages
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		title := strings.Join(r.Form["title"], "")
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	}
	f, err := db.GetTitles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "geting titles: %v\n", err)
	}
	if r.URL.Path[1:] == "frontpage" {
		renderFrontpage(w, r, f)
		return
	}
	renderTemplate(w, r, r.URL.Path[1:], &db.Page{})
}

// ViewHandler handles requests for viewing subpages
func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := db.PageByTitle(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, r, "view", p)
}

// EditHandler handles requests for editing subpages
func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := db.PageByTitle(title)
	if err != nil {
		p = &db.Page{Title: title}
		fmt.Fprintf(os.Stderr, "Page doesnt exist - created new Page obj: %v\n", err)
	}
	renderTemplate(w, r, "edit", p)
}

// SaveHandler handles requests for saving edited subpages.
// Checks if page exist, then update it, if not - inserts a new one
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &db.Page{Title: title, Body: []byte(body)}
	_, err := db.PageByTitle(title)
	if err == nil {
		err = db.UpdatePage(p)
	} else {
		_, err = db.InsertPage(p)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// MakeHandler extracts page title from the request and call a provided handler
func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/notfound", http.StatusFound)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, p *db.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		http.Redirect(w, r, "/notfound", http.StatusFound)
	}
}

func renderFrontpage(w http.ResponseWriter, r *http.Request, f *db.Frontpage) {
	err := templates.ExecuteTemplate(w, "frontpage.html", f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		http.Redirect(w, r, "/notfound", http.StatusFound)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", fmt.Errorf("invalid page title")
	}
	return m[2], nil
}

var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html", "template/frontpage.html", "template/notfound.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
