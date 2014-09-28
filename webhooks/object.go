package webhooks

import(
    "gopkg.in/mgo.v2/bson"
    "log"
)

type ObjectFormat struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Message string `bson:"message,omitempty"`
	Fields map[string]*ObjectFormat `bson:"fields,omitempty"`
	Placeholder string `bson:"placeholder,omitempty"`
}

func GetSampleObjectFormat() *ObjectFormat{
	return &ObjectFormat{
		Id : bson.NewObjectId(),
		Message : `New order worth $$$Subtotal_price$$`,
		Fields : map[string]*ObjectFormat{
			"subtotal_price" : &ObjectFormat{
				Message : "",
				Fields: nil,
				Placeholder: "$$Subtotal_price$$",
			},
			"line_items" : &ObjectFormat{
				Message : `- $$Quantity$$ $$Name$$`,
				Fields : map[string]*ObjectFormat{
					"quantity" : &ObjectFormat{
						Message : "",
						Fields: nil,
						Placeholder: "$$Quantity$$",
					},
					"name" : &ObjectFormat{
						Message : "",
						Fields: nil,
						Placeholder : "$$Name$$",
					},
				},
				Placeholder : "",
			},
		},
		Placeholder: "",
	}
}

func GetObjectFormat(ctx *Context, formatId string) *ObjectFormat{
	c := ctx.DB.C("object_format")
	res := new(ObjectFormat)
	err := c.FindId(bson.ObjectIdHex(formatId)).One(&res)
	if(err != nil){
		log.Printf("%+v", err)
	}
	return res
}

func (o *ObjectFormat) getMessage() string{
	return o.Message
}

func (o *ObjectFormat) FieldExists(fieldName string) bool{
	_, exists := o.Fields[fieldName]
	return exists
}

func (o *ObjectFormat) getFieldValue(fieldName string) string{
	val := o.Fields[fieldName].Placeholder
	return val
}

func (o *ObjectFormat) getChildObj(child string) *ObjectFormat{
	childObj := o.Fields[child]
	return childObj
}