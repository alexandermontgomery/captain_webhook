package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Transformer struct {
	Id                   bson.ObjectId                    `bson:"_id,omitempty" json:"Id,omitempty"`
	Name                 string                           `bson:"Name" json:"Name,omitempty"`
	ObjectTransformation map[string]*ObjectTransformation `bson:"ObjectTransformation" json:"ObjectTransformation"`
	WrapTemplate         string                           `bson:"WrapTemplate" json:"WrapTemplate"`
}

type ObjectTransformation struct {
	Rel_id   string `bson:"Rel_id" json:"rel_id"`
	Template string `bson:"Template" json:"template"`
	Weight   int    `bson:"Weight" json:"Weight"`
}

func LoadTransformer(ctx *Context, transformerId string) (*Transformer, error) {
	c := ctx.DB.C("transformer")
	transformer := new(Transformer)
	err := c.FindId(bson.ObjectIdHex(transformerId)).One(&transformer)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	if transformer.ObjectTransformation == nil {
		transformer.ObjectTransformation = make(map[string]*ObjectTransformation)
	}

	return transformer, nil
}

func ListTransformers(ctx *Context, limit int) []Transformer {
	c := ctx.DB.C("transformer")
	var transformer Transformer
	var transformerList []Transformer
	iter := c.Find(nil).Limit(limit).Iter()
	for iter.Next(&transformer) {
		if transformer.ObjectTransformation == nil {
			transformer.ObjectTransformation = make(map[string]*ObjectTransformation)
		}
		transformerList = append(transformerList, transformer)
	}

	if err := iter.Close(); err != nil {
		return transformerList
	}

	return transformerList
}
