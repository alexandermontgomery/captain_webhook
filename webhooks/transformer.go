package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Transformer struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string        `bson:"name" json:"name,omitempty"`
	ObjectId     bson.ObjectId `bson:"object_id" json:"object_id,omitempty"`
	ObjectFormat *ObjectFormat `bson:",omitempty" json:"object,omitempty"`
}

func LoadTransformer(ctx *Context, transformerId string) (*Transformer, error) {
	c := ctx.DB.C("transformer")
	transformer := new(Transformer)
	err := c.FindId(bson.ObjectIdHex(transformerId)).One(&transformer)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}
	transformer.ObjectFormat = GetObjectFormat(ctx, transformer.ObjectId.Hex())

	return transformer, nil
}

func ListTransformers(ctx *Context, limit int) []Transformer {
	c := ctx.DB.C("transformer")
	var transformer Transformer
	var transformerList []Transformer
	var objectList []bson.ObjectId
	iter := c.Find(nil).Limit(limit).Iter()
	for iter.Next(&transformer) {
		objectList = append(objectList, transformer.ObjectId)
		transformerList = append(transformerList, transformer)
	}

	c = ctx.DB.C("object_format")
	iter = c.Find(bson.M{"_id": bson.M{"$in": objectList}}).Iter()
	var objectFormat ObjectFormat
	objectMap := make(map[bson.ObjectId]ObjectFormat)
	for iter.Next(&objectFormat) {
		objectMap[objectFormat.Id] = objectFormat
	}

	log.Printf("%+v", objectMap)

	for index, transformer := range transformerList {
		tmpObjectFormat := objectMap[transformer.ObjectId]
		transformerList[index].ObjectFormat = &tmpObjectFormat
	}

	if err := iter.Close(); err != nil {
		return transformerList
	}

	return transformerList
}
