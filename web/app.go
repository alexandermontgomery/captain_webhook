package main

import (
	"fmt"
	"net/http"
	"captain_webhook/webhooks"
	"bytes"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

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
func index(w http.ResponseWriter, req *http.Request, ctx *webhooks.Context) (err error) {
	fmt.Fprintf(w, "Hello World!!")

	c := ctx.DB.C("object_format")
	obj := webhooks.GetSampleObjectFormat();
	err = c.Insert(obj)
    if(err != nil){
    	log.Printf("%+v", err)
    }
	return
}

func webhookHandle(w http.ResponseWriter, req *http.Request, ctx *webhooks.Context) (err error) {
	// params := req.URL.Query()
	// id := params.Get(":id")
	vars := mux.Vars(req)
	id := vars["id"]
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	obj := webhooks.GetObjectFormat(ctx, id)
	msg := webhooks.NewMessage(buf.Bytes(), obj)
    go webhooks.Publish(msg)
    return
}