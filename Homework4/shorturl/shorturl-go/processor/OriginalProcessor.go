package processor

import (
	"container/list"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type OriginalProcessor struct {
	*BaseProcessor
}

const POST string = "POST"
const TOKEN string = "token"
const ORIGINAL_URL string = "url"
const SHORT_URL string = "short"

func (this *OriginalProcessor) ProcessRequest(method, requestUrl string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error {
	if method != POST {
		return errors.New("Create short url must be POST")
	}

	var bodyInfo map[string]interface{}
	err := json.Unmarshal(body, &bodyInfo)
	if err != nil {
		return err
	}

	originalUrl, has := bodyInfo[ORIGINAL_URL].(string)
	if !has {
		return errors.New("Post info errors")
	}

	shortUrl, err := this.createUrl(originalUrl)
	if err != nil {
		return err
	}

	response, err := this.createResponseJson(shortUrl)
	if err != nil {
		return err
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")
	io.WriteString(w, response)

	return nil
}

func (this *OriginalProcessor) createUrl(originalUrl string) (string, error) {
	short, err := this.Lru.GetShortURL(originalUrl)
	if err == nil {
		return short, nil
	}

	count, err := this.CountFunction()
	if err != nil {
		return "", nil
	}

	shortUrl, err := transNumToString(count)
	if err != nil {
		return "", err
	}

	this.Lru.SetURL(originalUrl, shortUrl)
	return shortUrl, nil

}

func (this *OriginalProcessor) createResponseJson(shortUrl string) (string, error) {
	jsonResponse := make(map[string]interface{})
	jsonResponse[SHORT_URL] = shortUrl

	res, err := json.Marshal(jsonResponse)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func transNumToString(num int64) (string, error) {

	var base int64
	base = 62
	baseHex := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	output_list := list.New()
	for num/base != 0 {
		output_list.PushFront(num % base)
		num = num / base
	}
	output_list.PushFront(num % base)
	str := ""
	for iter := output_list.Front(); iter != nil; iter = iter.Next() {
		str = str + string(baseHex[int(iter.Value.(int64))])
	}

	return str, nil
}
