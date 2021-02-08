package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	Router *mux.Router
	Server *http.Server
	Store  Store
	Client Client
	Conn   *sql.DB
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

	a.Store = Store{}
	a.Client = Client{}
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
	a.Router.Path("/signup").HandlerFunc(a.Store.saveUser).Methods("POST")
	a.Router.Path("/login").HandlerFunc(a.Store.login).Methods("POST")
}

func initializeDatabaseConnection() {
	log.Println("Initializing connection to the database...")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	a.Conn, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
		return
	}
	log.Println("Connection to database has been established")
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
	defer a.Conn.Close()
	log.Println("Attempting to run the server...")
	go runServerPipe(serverDone)
	isServerDone := <-serverDone
	log.Println("Attempting to initialize connection to database...")
	initializeDatabaseConnection()
	log.Println(isServerDone)
	log.Println("Attempting to run the client")
	go a.Client.runClientPipe(clientDone)
	<-clientDone
}
