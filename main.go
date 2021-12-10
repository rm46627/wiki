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
	err := db.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// handling requests from client
	http.HandleFunc("/", wiki.Handler)
	http.HandleFunc("/frontpage", wiki.Handler)
	http.HandleFunc("/view/", wiki.MakeHandler(wiki.ViewHandler))
	http.HandleFunc("/edit/", wiki.MakeHandler(wiki.EditHandler))
	http.HandleFunc("/save/", wiki.MakeHandler(wiki.SaveHandler))
	http.HandleFunc("/delete/", wiki.MakeHandler(wiki.DeleteHandler))

	// files
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(":80", nil))
}
