package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"text/template"
	"time"
)

func ReceiveMessage(data []byte, transformer_id string) *Message {
	msg := new(Message)
	msg.Body = data
	msg.ParseAndFlatten()
	msg.Id = bson.NewObjectId()
	msg.TransformerId = bson.ObjectIdHex(transformer_id)
	msg.Time = time.Now()
	return msg
}

func (msg *Message) ParseAndFlatten() {
	msg.Json = parse(msg.Body)
	msg.BodyStr = string(msg.Body)
	msg.Flat = make(map[string]*MessageObject)
	FlattenMessage(msg.Json, msg.Flat, "")
}

func FlattenMessage(json map[string]interface{}, flat map[string]*MessageObject, rel_id string) {

	divider := "::"

	if rel_id == "" {
		rel_id = "root"
	}

	data := make(map[string]interface{})
	is_arr := false
	for k, v := range json {
		switch v.(type) {
		case string:
			data[k] = v
		case int:
			data[k] = v
		case nil:
			data[k] = nil
		case bool:
			data[k] = v
		case []interface{}:
			is_arr = true
		case map[string]interface{}:
			obj, _ := v.(map[string]interface{})
			FlattenMessage(obj, flat, rel_id+divider+k)
		}
	}

	msgObj := new(MessageObject)
	msgObj.Data = data
	msgObj.Rel_id = rel_id
	msgObj.Is_array = is_arr
	flat[rel_id] = msgObj
	return
}

func parse(data []byte) map[string]interface{} {
	var obj interface{}
	err := json.Unmarshal(data, &obj)
	newObj := obj.(map[string]interface{})
	if err != nil {
		log.Println(err)
		return newObj
	}
	return newObj
}

type Message struct {
	Id            bson.ObjectId             `bson:"_id,omitempty" json:"Id"`
	Json          map[string]interface{}    `json:"-"`
	Body          []byte                    `bson:"body,omitempty" json:"-"`
	BodyStr       string                    `bson:"-"`
	Flat          map[string]*MessageObject `json:"Flat"`
	Time          time.Time
	TransformerId bson.ObjectId `bson:"transformer_id" json:"TransformerId"`
	translatedStr string
}

type MessageObject struct {
	Data     map[string]interface{}
	Rel_id   string
	Is_array bool
}

func (msg *Message) Translate() {
	var translatedStr bytes.Buffer
	t, err := template.New("").Parse("msg.obj.getMessage()")
	if err != nil {
		log.Printf("Error translating message: %s", err)
	}

	execErr := t.Execute(&translatedStr, nil)
	if execErr != nil {
		log.Printf("Error translating message: %s", err)
	}
	msg.translatedStr = string(translatedStr.Bytes())
	return
}

func (msg *Message) LogMessage(ctx *Context) {
	c := ctx.DB.C("messages")
	err := c.Insert(msg)
	if err != nil {
		log.Printf("Problem inserting log Message: %+v", err)
	}
}

func LoadMessageLog(ctx *Context, transformerId string, limit int) ([]Message, error) {
	c := ctx.DB.C("messages")
	var message Message
	var err error
	messageList := make([]Message, 0)
	transformerIdObj := bson.ObjectIdHex(transformerId)
	iter := c.Find(bson.M{"transformer_id": transformerIdObj}).Limit(limit).Iter()

	for iter.Next(&message) {
		message.ParseAndFlatten()
		messageList = append(messageList, message)
	}

	if err = iter.Close(); err != nil {
		log.Printf("Problem loading log messages: %+v", err)
		return messageList, err
	}

	return messageList, err
}
