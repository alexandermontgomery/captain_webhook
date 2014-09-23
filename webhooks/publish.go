package webhooks

import (
	"log"
	"net/http"
	"net/url"
	"bytes"
	"fmt"
	"io/ioutil"
)

func Publish(msg *Message){
	hookUrl := "https://outdoorvoices.slack.com/services/hooks/incoming-webhook?token=ycyxBoayr39SgsdqMmsTlUBK	"
	jsonStr := `{"text":"` + msg.translatedStr + `"}`
	log.Println(jsonStr)

	data := url.Values{}
	data.Set("payload", jsonStr)

	req, err := http.NewRequest("POST", hookUrl, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

    resp, err := client.Do(req)

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    log.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    resp_bod, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(resp_bod))
}