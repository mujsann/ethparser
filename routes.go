package main

import "net/http"

func defineRoutes(app *App) {
	http.HandleFunc("/current-block", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			app.getCurrentBlockHandler(w, r)
		} else {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
		}
	})
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			app.subscribeHandler(w, r)
		} else {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
		}
	})
	http.HandleFunc("/transactions/{address}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			app.getTransactionsHandler(w, r)
		} else {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
		}
	})
}
