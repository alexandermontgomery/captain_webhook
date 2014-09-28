package main

import (
	"fmt"
	"net/http"
	"captain_webhook/webhooks"
	"bytes"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

// declare a new key type for context to avoid collisions 
type ctx_key int
const CTX_DB ctx_key = 0

var dbsession *mgo.Session

func main() {
	var err error
	dbsession, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(fmt.Sprintf("Can't connect to mongo, go error %v\n", err))
	} else {
		log.Println("Connected to Mongo")
	}
	cwd, _ := os.Getwd()
	log.Printf("Current Working Directory: %s", cwd)
	router := mux.NewRouter()
	router.Handle("/", MiddlewareHandler(index)).Methods("GET")
  	router.Handle("/webhook/{id}", MiddlewareHandler(webhookHandle)).Methods("POST")
  	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/vagrant/src/static/")))).Methods("GET")

	http.Handle("/", router)
	addr := ":8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type MiddlewareHandler func(http.ResponseWriter, *http.Request) error

func (mh MiddlewareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	conn := dbsession.Copy().DB("captain_webhook")
	defer conn.Session.Close()
	
	context.Set(req, CTX_DB, conn)

	//run the handler and grab the error, and report it
	err := mh(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDBConn(req *http.Request) (*mgo.Database){
	return context.Get(req, CTX_DB).(*mgo.Database)
}

// PAGE CALLBACKS
func index(w http.ResponseWriter, req *http.Request) (err error) {
	fmt.Fprintf(w, "Hello World!!")

	db_conn := getDBConn(req)
	c := db_conn.C("object_format")
	obj := webhooks.GetSampleObjectFormat();
	err = c.Insert(obj)
    if(err != nil){
    	log.Printf("%+v", err)
    }
	return
}

func webhookHandle(w http.ResponseWriter, req *http.Request) (err error) {
	// params := req.URL.Query()
	// id := params.Get(":id")
	vars := mux.Vars(req)
	id := vars["id"]
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	dbConn := getDBConn(req)
	obj := webhooks.GetObjectFormat(dbConn, id)
	msg := webhooks.NewMessage(buf.Bytes(), obj)
    go webhooks.Publish(msg)
    return
}