package processor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	config "../config"
)

type Router struct {
	Configure *config.Configure
	Processor map[int]Processor
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var action int
	if r.Method == "GET" {
		action = 0
	} else {
		action = 1
	}
	requestUrl := r.RequestURI[1:]
	err := r.ParseForm()
	if err != nil {
		return
	}
	params := make(map[string]string)
	for k, v := range r.Form {
		params[k] = v[0]
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		return
	}

	processor, _ := this.Processor[action]
	err = processor.ProcessRequest(r.Method, requestUrl, params, body, w, r)
	if err != nil {
		fmt.Printf("[ERROR] : %v\n", err)
	}
	return

}
