package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/smth/view", app.smthView)
	mux.HandleFunc("/smth/viewbyid", app.smthViewById)
	mux.HandleFunc("/smth/create", app.smthCreate)
	return mux
}
