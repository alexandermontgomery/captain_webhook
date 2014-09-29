package webhooks

import(
    "gopkg.in/mgo.v2/bson"
    "log"
)

type Transformer struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string `bson:"name"`
	ObjectId bson.ObjectId `bson:"object_id"`
	ObjectFormat *ObjectFormat `bson:",omitempty"`
}

func LoadTransformer(ctx *Context, transformerId string) (*Transformer, error){
	c := ctx.DB.C("transformer")
	transformer := new(Transformer)
	err := c.FindId(bson.ObjectIdHex(transformerId)).One(&transformer)
	if(err != nil){
		log.Printf("%+v", err)
		return nil, err
	}
	log.Printf("%+v", transformer)
	transformer.ObjectFormat = GetObjectFormat(ctx, transformer.ObjectId.Hex())

	return transformer, nil
}