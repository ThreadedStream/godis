package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	Server *http.Server
	Store  *Store
}

func (a *App) initializeApp(addr string) {
	a.Router = mux.NewRouter()

	a.Server = &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a.Store = &Store{}
	a.Store.initStore()
	a.initRoutes()
}

func (a *App) initRoutes() {
	a.Router.Path("/set").HandlerFunc(a.Store.SetKey).Methods("POST")
	a.Router.Path("/hset").HandlerFunc(a.Store.HSet).Methods("POST")
	a.Router.Path("/get").HandlerFunc(a.Store.GetKey).Methods("GET")
	a.Router.Path("/keys").HandlerFunc(a.Store.Keys).Methods("GET")
	a.Router.Path("/del").HandlerFunc(a.Store.Del).Methods("DELETE")
}

func (a *App) runStuff() {
	log.Fatal(a.Server.ListenAndServe())
}

func main() {
	addr := "0.0.0.0:5680"
	a := App{}
	a.initializeApp(addr)
	a.initRoutes()
	log.Println("Running server on " + addr)
	a.runStuff()
}
