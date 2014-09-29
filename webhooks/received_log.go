package webhooks

import(
	"log"
	"gopkg.in/mgo.v2/bson"
)

type MessageLog struct {
	Id bson.ObjectId `bson:"_id"`
	Body string `bson:"body"`
	Time  bson.MongoTimestamp `bson:"time"`
	ConfigId bson.ObjectId `bson:"config_id"`
}

func LogMessage(ctx *Context, data []byte, configId string){
	msgLog := new(MessageLog)
	msgLog.Id = bson.NewObjectId()
	msgLog.ConfigId = bson.ObjectIdHex(configId)
	msgLog.Body = string(data)

	c := ctx.DB.C("message_log")
	err := c.Insert(msgLog)
	if err != nil{
		log.Printf("Problem inserting log Message: %+v", err)
	}
}