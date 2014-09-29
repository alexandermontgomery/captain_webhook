package main

import (
	"fmt"
	"net/http"
	"captain_webhook/webhooks"
	"bytes"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

var dbsession *mgo.Session
var homeDir string = "/vagrant/src/src/captain_webhook"
var view *View

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
  	router.Handle("/configs/list", MiddlewareHandler(configListHandle)).Methods("GET")
  	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(homeDir + "/static/")))).Methods("GET")

	http.Handle("/", router)
	addr := ":8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func ensureIndices(conn *mgo.Database) (err error){
	index := mgo.Index{
	    Key: []string{"transformer_id"},
	    Unique: false,
	    DropDups: false,
	    Background: true, // See notes.
	    Sparse: false,
	}
	err = conn.C("message_log").EnsureIndex(index)
	if err != nil{
		return err
	}
	return err
}

type MiddlewareHandler func(http.ResponseWriter, *http.Request, *webhooks.Context) error

func (mh MiddlewareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  //create the context
  ctx, err := webhooks.NewContext(req, dbsession)
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

// PAGE CALLBACKS
func indexHandle(w http.ResponseWriter, req *http.Request, ctx *webhooks.Context) (err error) {
	view.RenderPage(w, "home", map[string]interface{}{"a" : 123, "b" : "Test String"})
	var a webhooks.Transformer
	a = webhooks.Transformer{
		bson.NewObjectId(),
		"Shopify",
		bson.ObjectIdHex("54274f086807a32854000001"),
		nil,
	}
	ctx.DB.C("transformer").Insert(a)
	return
}

func configListHandle(w http.ResponseWriter, req *http.Request, ctx *webhooks.Context) (err error) {
	view.RenderPage(w, "webhook_edit", nil)
	return
}

func webhookHandle(w http.ResponseWriter, req *http.Request, ctx *webhooks.Context) (err error) {
	// params := req.URL.Query()
	// id := params.Get(":id")
	vars := mux.Vars(req)
	id := vars["id"]
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	transformer, err := webhooks.LoadTransformer(ctx, id)
	if err != nil {
		return err
	}
	obj := transformer.ObjectFormat
	webhooks.LogMessage(ctx, buf.Bytes(), id)
	msg := webhooks.NewMessage(buf.Bytes(), obj)
    go webhooks.Publish(msg)
    return
}