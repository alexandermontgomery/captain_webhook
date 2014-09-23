package webhooks

type Object struct {
	fields map[string]interface{}
	message string
}

func NewObject(webhook_type string) *Object{
	newFields := new(Object)
	newFields.LoadObject(webhook_type)
	return newFields
}

func (o *Object) LoadObject(webhook_type string){
	o.message = `New order worth $$$Subtotal_price$$`
	o.fields = map[string]interface{}{
		"line_items" : Object{
			map[string]interface{}{
				"quantity" : "$$Quantity$$",
				"name" : "$$Name$$",
			},
			`- $$Quantity$$ $$Name$$`,
		},
		"subtotal_price" : "$$Subtotal_price$$",
	}
}

func (o *Object) getMessage() string{
	return o.message
}

func (o *Object) FieldExists(fieldName string) bool{
	_, exists := o.fields[fieldName]
	return exists
}

func (o *Object) getFieldValue(fieldName string) string{
	val, _ := o.fields[fieldName]
	return val.(string)
}

func (o *Object) getChildObj(child string) *Object{
	childObj := o.fields[child].(Object)
	return &childObj
}