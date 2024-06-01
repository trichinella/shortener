package main

import "net/http"

func start() {
	mux := http.NewServeMux()
	port := "8080"

	s := CreateStore("http://localhost:", "8080")

	mux.HandleFunc(`/`, mainPage(s))
	err := http.ListenAndServe(`:`+port, mux)
	if err != nil {
		panic(err)
	}
}
