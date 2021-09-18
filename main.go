package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rm46627/wiki/db"
	"github.com/rm46627/wiki/wiki"
)

func main() {

	// initializing data base
	err := db.Init("data.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// handling requests from client
	http.HandleFunc("/view/", wiki.MakeHandler(wiki.ViewHandler))
	http.HandleFunc("/edit/", wiki.MakeHandler(wiki.EditHandler))
	http.HandleFunc("/save/", wiki.MakeHandler(wiki.SaveHandler))

	// files
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
