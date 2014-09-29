package webhooks

import (
	"encoding/json"
	"log"
	"strings"
    "text/template"
    "bytes"
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

func NewMessage(data []byte, obj *ObjectFormat) *Message{
	msg := new(Message)
	jsonObj := parse(data)
	msg.jsonObj = jsonObj.(map[string]interface{})
	msg.obj = obj
	msg.Translate()
	return msg
}

func NewChildMessage(parent *Message, data interface{}, childObj *ObjectFormat) *Message{
	child := new(Message)
	child.jsonObj = data.(map[string]interface{})
	child.obj = childObj
	child.Translate()
	parent.children = append(parent.children, child)
	child.parent = parent
	return child
}

type Message struct {
	jsonObj map[string]interface{}
	translatedStr string
	obj *ObjectFormat
	children []*Message
	parent *Message
}

func (msg *Message) Translate() {
    vars := make(map[string]interface{})    
    for k,v := range msg.jsonObj {
    	inter, isInter := v.([]interface{})

    	if(!msg.obj.FieldExists(k)){
    		continue
    	} 

    	if(isInter){
    		childObj := msg.obj.getChildObj(k)
            var strArr []string
    		for _, u := range inter {
	    		child := NewChildMessage(msg, u, childObj)
    			strArr = append(strArr, child.translatedStr)
	        }    
            log.Printf("%+v", strArr)        
            vars[msg.obj.getFieldPlaceholder(k)] = strings.Join(strArr, `\n`)

    	} else{
            vars[msg.obj.getFieldPlaceholder(k)] = v
        }    	
    }

    var translatedStr bytes.Buffer

    t, err := template.New("").Parse(msg.obj.getMessage())
    if err != nil {
        log.Printf("Error translating message: %s", err)
    }
    execErr := t.Execute(&translatedStr, vars);
    if execErr != nil {
        log.Printf("Error translating message: %s", err)
    }
    msg.translatedStr = string(translatedStr.Bytes())
}