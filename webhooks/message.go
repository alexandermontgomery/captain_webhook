package webhooks

import (
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

func parse(data []byte) interface{} {
	var obj interface{}
	err := json.Unmarshal(data, &obj)
    if err != nil {
        log.Println(err)
        return obj
    }

    return obj
}

func NewMessage(webhook_type string, data []byte) *Message{
	msg := new(Message)
	jsonObj := parse(data)
	msg.jsonObj = jsonObj.(map[string]interface{})
	msg.webhookType = webhook_type
	msg.obj = NewObject(webhook_type)
	msg.Translate()
	return msg
}

func NewChildMessage(parent *Message, data interface{}, childObj *Object) *Message{
	child := new(Message)
	child.jsonObj = data.(map[string]interface{})
	child.webhookType = parent.webhookType
	child.obj = childObj
	child.Translate()
	parent.children = append(parent.children, child)
	child.parent = parent
	return child
}

type Message struct {
	jsonObj map[string]interface{}
	webhookType string
	translatedStr string
	obj *Object
	children []*Message
	parent *Message
}

func (msg *Message) Translate() {
	msg.translatedStr = msg.obj.getMessage()
    for k,v := range msg.jsonObj {
    	inter, isInter := v.([]interface{})

    	if(!msg.obj.FieldExists(k)){
    		continue
    	} 

    	if(isInter){
    		childObj := msg.obj.getChildObj(k)
    		for _, u := range inter {
	    		child := NewChildMessage(msg, u, childObj)
    			s := []string{msg.translatedStr, child.translatedStr}
    			msg.translatedStr = strings.Join(s, `\n`)
	        }
    	}
    	var strValue string

    	switch vv := v.(type){
    		case string:
    			strValue = vv    			
    		case int:
    			strValue = strconv.Itoa(v.(int))
    		case float64:
    			strValue = strconv.FormatFloat(vv, 'f', 2, 64)
    		case bool:
    			strValue = strconv.FormatBool(vv)

    	}
    	if(len(strValue) > 0){
    		msg.translatedStr = strings.Replace(msg.translatedStr, msg.obj.getFieldValue(k), strValue, -1)   	
    	}
    }
}