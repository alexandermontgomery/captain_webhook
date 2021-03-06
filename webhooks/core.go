package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
)

type Context struct {
	DB *mgo.Database
}

func (ctx *Context) Close() {
	ctx.DB.Session.Close()
}

func NewContext(req *http.Request, dbsession *mgo.Session) (*Context, error) {
	ctx := new(Context)
	ctx.DB = dbsession.Copy().DB("captain_webhook")
	return ctx, nil
}
