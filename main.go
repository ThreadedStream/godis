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
	Client *Client
}

//For simplification
var a *App = &App{}

func (a *App) initializeApp(addr string) {
	a.Router = mux.NewRouter()

	a.Server = &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a.Store = &Store{}
	a.Client = &Client{}
	a.Store.initStore()
	a.Client.initParserPipe()
	a.initRoutes()
}

func (a *App) initRoutes() {
	a.Router.Path("/set").HandlerFunc(a.Store.Set).Methods("POST")
	a.Router.Path("/hset").HandlerFunc(a.Store.HSet).Methods("POST")
	a.Router.Path("/mset").HandlerFunc(a.Store.MSet).Methods("POST")
	a.Router.Path("/get").HandlerFunc(a.Store.Get).Methods("GET")
	a.Router.Path("/hget").HandlerFunc(a.Store.HGet).Methods("GET")
	a.Router.Path("/mget").HandlerFunc(a.Store.MGet).Methods("GET")
	a.Router.Path("/keys").HandlerFunc(a.Store.Keys).Methods("GET")
	a.Router.Path("/del").HandlerFunc(a.Store.Del).Methods("DELETE")
	a.Router.Path("/save").HandlerFunc(a.Store.SerializedStore).Methods("GET")
	a.Router.Path("/restore").HandlerFunc(a.Store.Restore).Methods("POST")
}

func (a *App) serverListenAndServe() {
	log.Fatal(a.Server.ListenAndServe())
}

func runServerPipe(done chan bool) {
	addr := "0.0.0.0:5680"
	a.initializeApp(addr)
	a.initRoutes()
	log.Println("Running server on " + addr)
	done <- true
	a.serverListenAndServe()
}

func main() {
	serverDone := make(chan bool)
	clientDone := make(chan bool)
	log.Println("Attempting to run the server...")
	go runServerPipe(serverDone)
	isServerDone := <-serverDone
	log.Println(isServerDone)
	log.Println("Attempting to run the client")
	go a.Client.runClientPipe(clientDone)
	<-clientDone
}
