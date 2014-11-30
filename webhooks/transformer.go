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

const DESTINATION_CONTENT_TYPE_JSON int = 0
const DESTINATION_CONTENT_TYPE_FORM_DATA int = 1

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
	Destination          TransformerDestination  `bson:"Destination" json:"Destination"`
	Active               bool                    `bson:"Active" json:"Active"`
}

type TransformerDestination struct {
	Url         string `bson:"Url" json:"Url"`
	ContentType int    `bson:"ContentType" json:"ContentType"`
	FormField   string `bson:"FormField" json:"FormField"`
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

	var reqData *bytes.Buffer
	var contentType string
	switch trans.Destination.ContentType {
	case DESTINATION_CONTENT_TYPE_JSON:
		reqData = bytes.NewBufferString(string(transStr))
		contentType = "application/json"
	case DESTINATION_CONTENT_TYPE_FORM_DATA:
		data := url.Values{}
		data.Set(trans.Destination.FormField, string(transStr))
		reqData = bytes.NewBufferString(data.Encode())
		contentType = "application/x-www-form-urlencoded"
	default:
		log.Printf("No Content-Type selected, request could not be made")
		return
	}

	req, err := http.NewRequest("POST", trans.Destination.Url, reqData)
	req.Header.Add("Content-Type", contentType)
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

	if !trans.Id.Valid() {
		trans.Id = bson.NewObjectId()
	}
	change, err := c.UpsertId(trans.Id, trans)

	if err != nil {
		log.Printf("Problems saving transformer: %+v", err)
	}
	return change, err
}

func DeleteTransformer(ctx *Context, transformerId string) {
	c := ctx.DB.C("transformer")
	err := c.RemoveId(bson.ObjectIdHex(transformerId))

	if err != nil {
		log.Printf("Problems saving transformer: %+v", err)
	}
	return
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
