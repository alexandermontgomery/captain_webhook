package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

var dbsession *mgo.Session
var homeDir string = "/vagrant/src/src/captain_webhook"
var view *View
var env = "dev"

func main() {
	var err error
	dbsession, err = mgo.Dial("localhost:27017")
	view = NewView()

	if err != nil {
		panic(fmt.Sprintf("Can't connect to mongo, go error %v\n", err))
	} else {
		log.Println("Connected to Mongo")
	}

	indicesErr := ensureIndices(dbsession.DB("captain_webhook"))
	if indicesErr != nil {
		panic(fmt.Sprintf("Could not ensure all indices %+v\n", err))
	}

	cwd, _ := os.Getwd()
	log.Printf("Current Working Directory: %s", cwd)
	router := mux.NewRouter()
	router.Handle("/", MiddlewareHandler(indexHandle)).Methods("GET")
	router.Handle("/webhook/{id}", MiddlewareHandler(webhookHandle)).Methods("POST")
	router.Handle("/transformers", MiddlewareHandler(transformersListHandle)).Methods("GET")
	router.Handle("/transformers/{id}", MiddlewareHandler(transformersSingleHandle)).Methods("GET")
	router.Handle("/transformers/{id}", MiddlewareHandler(transformersSaveHandle)).Methods("PUT")
	router.Handle("/transformers/{id}/messages", MiddlewareHandler(transformersMessagesHandle)).Methods("GET")
	router.Handle("/templates/{name}", http.HandlerFunc(templateServerHandle)).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(homeDir+"/static/")))).Methods("GET")

	http.Handle("/", router)
	addr := ":8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func ensureIndices(conn *mgo.Database) (err error) {
	index := mgo.Index{
		Key:        []string{"transformer_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     false,
	}
	err = conn.C("message_log").EnsureIndex(index)
	if err != nil {
		return err
	}
	return err
}

type MiddlewareHandler func(http.ResponseWriter, *http.Request, *Context) error

func (mh MiddlewareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := NewContext(req, dbsession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer ctx.Close()

	//run the handler and grab the error, and report it
	err = mh(w, req, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func templateServerHandle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	view.ServeSingleTemplate(w, name)
	return
}

func writeJson(w http.ResponseWriter, b []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

// PAGE CALLBACKS
func indexHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	view.RenderPage(w, "home", nil)
	return
}

func configListHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	view.RenderPage(w, "webhook_edit", nil)
	return
}

func transformersListHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	transformers := ListTransformers(ctx, 10)
	// enc := json.NewEncoder(w)
	// enc.Encode(transformers)
	out, err := json.Marshal(transformers)
	if err != nil {
		return
	}
	writeJson(w, out, 200)
	return
}

func transformersSingleHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	vars := mux.Vars(req)
	id := vars["id"]

	transformer, err := LoadTransformer(ctx, id)
	out, err := json.Marshal(transformer)
	if err != nil {
		return
	}
	writeJson(w, out, 200)
	return
}

func transformersMessagesHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	vars := mux.Vars(req)
	id := vars["id"]

	messagesLog, err := LoadMessageLog(ctx, id, 10)
	if err != nil {
		return
	}

	out, err := json.Marshal(messagesLog)
	if err != nil {
		return
	}
	writeJson(w, out, 200)
	return
}

func transformersSaveHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	var trans Transformer

	err = json.Unmarshal(buf.Bytes(), &trans)

	if err != nil {
		log.Printf("Could not unmarshal Json on transformer save: %+v", err)
		return
	}

	_, err = SaveTransformer(ctx, &trans)

	return
}

func webhookHandle(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	vars := mux.Vars(req)
	id := vars["id"]
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	msg := ReceiveMessage(buf.Bytes(), id)
	msg.LogMessage(ctx)
	out, err := json.Marshal(msg)
	writeJson(w, out, 200)
	//go Publish(msg)
	return
}
