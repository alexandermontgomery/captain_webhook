package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"text/template"
	"time"
)

func ReceiveMessage(ctx *Context, data []byte, transformer_id string) *Message {
	msg := new(Message)
	msg.Body = data
	msg.Id = bson.NewObjectId()
	msg.TransformerId = bson.ObjectIdHex(transformer_id)
	msg.Time = time.Now()
	msg.LogMessage(ctx)
	msg.ParseAndFlatten()
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
		case []interface{}: // If this is an array determine if it is an array of primitives or objects that we should flatten
			arr, _ := v.([]interface{})
			for _, w := range arr {
				switch w.(type) {
				case map[string]interface{}:
					part, _ := w.(map[string]interface{})
					FlattenMessage(part, flat, rel_id+divider+k)
				}
			}
		case map[string]interface{}:
			obj, _ := v.(map[string]interface{})
			FlattenMessage(obj, flat, rel_id+divider+k)
		}
	}

	if flat[rel_id] == nil {
		msgObj := new(MessageObject)
		msgObj.Data = data
		msgObj.Rel_id = rel_id
		flat[rel_id] = msgObj
	} else {
		flat[rel_id].Array_data = append(flat[rel_id].Array_data, data)
	}

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
	Data       map[string]interface{}
	Array_data []map[string]interface{}
	Rel_id     string
}

func (msg *Message) LogMessage(ctx *Context) {
	c := ctx.DB.C("messages")
	err := c.Insert(msg)
	if err != nil {
		log.Printf("Problem inserting log Message: %+v", err)
	}
}

func (msg *Message) Transform(trans *Transformer) ([]byte, error) {
	var transformation []byte
	var err error
	for _, v := range trans.ObjectTransformation {
		if msg.Flat[v.Rel_id] == nil {
			continue
		}

		t, err := template.New("").Parse(v.Template)
		if err != nil {
			log.Printf("Error transforming object: %+v", err)
			return transformation, err
		}

		var transBuf bytes.Buffer
		err = t.Execute(&transBuf, msg.Flat[v.Rel_id].Data)
		if err != nil {
			log.Printf("Error transforming object during template execution: %+v", err)
			return transformation, err
		}

		if len(msg.Flat[v.Rel_id].Array_data) > 0 {
			for _, vv := range msg.Flat[v.Rel_id].Array_data {
				transBuf.Write([]byte(", "))
				err = t.Execute(&transBuf, vv)
				if err != nil {
					log.Printf("Error transforming object during template execution: %+v", err)
					return transformation, err
				}
			}
		}

		bytesArr := transBuf.Bytes()

		if len(transformation) > 0 {
			bytesArr = append([]byte(" "), bytesArr...)
		}
		transformation = append(transformation, bytesArr...)
	}

	tplPrefix := []byte(trans.TemplatePrefix)
	tplSuffix := []byte(trans.TemplateSuffix)
	if len(tplPrefix) > 0 {
		transformation = append(tplPrefix, transformation...)
	}

	if len(tplSuffix) > 0 {
		transformation = append(transformation, tplSuffix...)
	}

	return transformation, err
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

func LoadMessage(ctx *Context, messageId string) (*Message, error) {
	c := ctx.DB.C("messages")
	var message *Message
	err := c.FindId(bson.ObjectIdHex(messageId)).One(&message)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	message.ParseAndFlatten()
	return message, err
}
