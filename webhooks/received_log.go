package webhooks

import(
	"log"
	"gopkg.in/mgo.v2/bson"
)

type MessageLog struct {
	Id bson.ObjectId `bson:"_id"`
	Body string `bson:"body"`
	Time  bson.MongoTimestamp `bson:"time"`
	TransformerId bson.ObjectId `bson:"transformer_id"`
}

func LogMessage(ctx *Context, data []byte, transformerId string){
	msgLog := new(MessageLog)
	msgLog.Id = bson.NewObjectId()
	msgLog.TransformerId = bson.ObjectIdHex(transformerId)
	msgLog.Body = string(data)

	c := ctx.DB.C("message_log")
	err := c.Insert(msgLog)
	if err != nil{
		log.Printf("Problem inserting log Message: %+v", err)
	}
}