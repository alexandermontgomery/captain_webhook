package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type MessageLog struct {
	Id            bson.ObjectId       `bson:"_id"`
	Body          string              `bson:"body"`
	Time          bson.MongoTimestamp `bson:"time"`
	TransformerId bson.ObjectId       `bson:"transformer_id"`
}

func LogMessage(ctx *Context, data []byte, transformerId string) {
	msgLog := new(MessageLog)
	msgLog.Id = bson.NewObjectId()
	msgLog.TransformerId = bson.ObjectIdHex(transformerId)
	msgLog.Body = string(data)

	c := ctx.DB.C("message_log")
	err := c.Insert(msgLog)
	if err != nil {
		log.Printf("Problem inserting log Message: %+v", err)
	}
}

func LoadMessageLog(ctx *Context, transformerId string, limit int) (messageLogList []MessageLog, err error) {
	c := ctx.DB.C("message_log")
	var messageLog MessageLog
	transformerIdObj := bson.ObjectIdHex(transformerId)
	iter := c.Find(bson.M{"transformer_id": transformerIdObj}).Limit(limit).Iter()

	for iter.Next(&messageLog) {
		messageLogList = append(messageLogList, messageLog)
	}

	if err = iter.Close(); err != nil {
		log.Printf("Problem loading log messages: %+v", err)
		return
	}

	return
}
