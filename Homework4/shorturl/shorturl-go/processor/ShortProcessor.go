package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ShortProcessor struct {
	*BaseProcessor
}

func (this *ShortProcessor) ProcessRequest(method, requestUrl string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error {
	originalUrl, err := this.GetOriginalUrl(requestUrl)
	if err != nil {
		response, err := this.createErrJson("not legal url")
		if err != nil {
			return err
		}
	
		header := w.Header()
		header.Add("Content-Type", "application/json")
		header.Add("charset", "UTF-8")
		io.WriteString(w, response)
	
		return nil
	}

	fmt.Printf("REQUEST_URL: %v --- ORIGINAL_URL : %v \n", requestUrl, originalUrl)
	response, err := this.createResponseJson(originalUrl)
	if err != nil {
		return err
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")
	io.WriteString(w, response)

	return nil
}

func (this *ShortProcessor) GetOriginalUrl(requestUrl string) (string, error) {
	originalUrl, err := this.Lru.GetOriginalURL(requestUrl)
	if err != nil {
		return "", err
	}

	return originalUrl, nil
}

func (this *ShortProcessor) createResponseJson(originalUrl string) (string, error) {
	jsonResponse := make(map[string]interface{})
	jsonResponse[ORIGINAL_URL] = originalUrl

	res, err := json.Marshal(jsonResponse)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (this *ShortProcessor) createErrJson(errMessage string) (string, error) {
	jsonResponse := make(map[string]interface{})
	jsonResponse["err"] = errMessage

	res, err := json.Marshal(jsonResponse)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
