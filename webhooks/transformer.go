package main

import (
	"bytes"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
)

type ByWeight []*ObjectTransformation

func (a ByWeight) Len() int {
	return len(a)
}
func (a ByWeight) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByWeight) Less(i, j int) bool {
	return a[i].Weight < a[j].Weight
}

type Transformer struct {
	Id                   bson.ObjectId           `bson:"_id,omitempty" json:"Id,omitempty"`
	Name                 string                  `bson:"Name" json:"Name,omitempty"`
	ObjectTransformation []*ObjectTransformation `bson:"ObjectTransformation" json:"ObjectTransformation"`
	TemplatePrefix       string                  `bson:"TemplatePrefix" json:"TemplatePrefix"`
	TemplateSuffix       string                  `bson:"TemplateSuffix" json:"TemplateSuffix"`
	DestinationUrl       string                  `bson:"DestinationUrl" json:"DestinationUrl"`
}

type ObjectTransformation struct {
	Rel_id   string `bson:"Rel_id" json:"Rel_id"`
	Template string `bson:"Template" json:"Template"`
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
		transformer.ObjectTransformation = make([]*ObjectTransformation, 0)
	}
	transformer.SortObjectTransformations()
	return transformer, nil
}

func (trans *Transformer) SortObjectTransformations() {
	sort.Sort(ByWeight(trans.ObjectTransformation))
}

func (trans *Transformer) PublishMessage(msg *Message) {
	transStr, _ := msg.Transform(trans)
	log.Printf("%+v", string(transStr))
	data := url.Values{}
	data.Set("payload", string(transStr))

	req, err := http.NewRequest("POST", trans.DestinationUrl, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	resp_bod, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(resp_bod))
}

func SaveTransformer(ctx *Context, trans *Transformer) (*mgo.ChangeInfo, error) {
	c := ctx.DB.C("transformer")

	var err error
	change, err := c.UpsertId(trans.Id, trans)

	if err != nil {
		log.Printf("Problems saving transformer: %+v", err)
	}
	return change, err
}

func ListTransformers(ctx *Context, limit int) []Transformer {
	c := ctx.DB.C("transformer")
	var transformer Transformer
	var transformerList []Transformer
	iter := c.Find(nil).Limit(limit).Iter()
	for iter.Next(&transformer) {
		if transformer.ObjectTransformation == nil {
			transformer.ObjectTransformation = make([]*ObjectTransformation, 0)
		}
		transformer.SortObjectTransformations()
		transformerList = append(transformerList, transformer)
	}

	if err := iter.Close(); err != nil {
		return transformerList
	}

	return transformerList
}
