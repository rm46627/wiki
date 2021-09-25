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

var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html", "template/frontpage.html", "template/notfound.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9~]+)$")

// Handler handles requests for front and notfound pages
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		title := strings.Join(r.Form["title"], "")
		http.Redirect(w, r, "/view/"+titleToURL(title), http.StatusFound) // zamiana title na string
	}
	f, err := db.GetPages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "geting pages: %v\n", err)
	}
	if r.URL.Path[1:] == "frontpage" {
		renderFrontpage(w, r, f)
		return
	}
	renderTemplate(w, r, r.URL.Path[1:], &db.Page{})
}

// ViewHandler handles requests for viewing subpages
func ViewHandler(w http.ResponseWriter, r *http.Request, url string) {
	p, err := db.PageByURL(url)
	if err != nil {
		http.Redirect(w, r, "/edit/"+url, http.StatusFound)
		return
	}
	renderTemplate(w, r, "view", p)
}

// EditHandler handles requests for editing subpages
func EditHandler(w http.ResponseWriter, r *http.Request, url string) {
	p, err := db.PageByURL(url)
	if err != nil {
		p = &db.Page{URL: url, Title: URLtoTitle(url)}
		fmt.Fprintf(os.Stderr, "Page does not exist - created new Page obj: %v\n", err)
	}
	renderTemplate(w, r, "edit", p)
}

// SaveHandler handles requests for saving edited subpages.
// Checks if page exist, then update it, if not - inserts a new one
func SaveHandler(w http.ResponseWriter, r *http.Request, url string) {
	body := r.FormValue("body")
	p := &db.Page{URL: url, Title: URLtoTitle(url), Body: []byte(body)}
	_, err := db.PageByURL(url)
	if err == nil {
		err = db.UpdatePage(p)
	} else {
		_, err = db.InsertPage(p)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	http.Redirect(w, r, "/view/"+url, http.StatusFound)
}

// DeleteHandler handles request for deleting subpages from the database
func DeleteHandler(w http.ResponseWriter, r *http.Request, url string) {
	err := db.DeletePage(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "deleting %s page: %v\n", url, err)
	}
	http.Redirect(w, r, "/frontpage", http.StatusFound)
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

// URLtoTitle convert url of subpage to its title.
// It turns '~' from url to space character.
// '~' still can by used in title but it must have spaces around it.
func URLtoTitle(s string) (title string) {
	var anotherTilde bool
	for _, x := range s {
		if x == '~' && anotherTilde == false {
			title += " "
			anotherTilde = true
		} else {
			title += string(x)
			anotherTilde = false
		}
	}
	return
}

// Works the opposite of urlToTitle func.
func titleToURL(title string) (s string) {
	var anotherTilde bool
	for _, x := range title {
		if x == ' ' && anotherTilde == false {
			s += "~"
			anotherTilde = true
		} else {
			s += string(x)
			anotherTilde = false
		}
	}
	return
}
