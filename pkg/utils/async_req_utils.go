package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-brrrr/pkg/types"
	"log"
	"net/http"
)

func setHeaders(req *http.Request, headers *map[string]string) {
	for key, value := range *headers {
		req.Header.Set(key, value)
	}
}

func SendRequest(
	url string,
	headers map[string]string,
	method string,
	postBody *map[string]interface{},
	response chan<- types.RespHolder,
) {
	// log.Printf("Sending request to %s", url)
	client := &http.Client{}
	var req *http.Request
	var clientResp *http.Response
	var err error
	res := types.RespHolder{Err: fmt.Errorf("Unable to communicate to %s", url)}
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			res.Err = err
			response <- res
			return
		}
	case http.MethodPost:
		var pb []byte
		pb, err := json.Marshal(postBody)
		if err != nil {
			res.Err = err
			response <- res
			return
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(pb))
	default:
		res.Err = fmt.Errorf("Method not supported: %s", method)
		response <- res
		return
	}
	setHeaders(req, &headers)
	clientResp, err = client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s", url)
		res.Err = err
		response <- res
		return
	}
	if clientResp.StatusCode > 399 {
		log.Printf("Client <%s> Responded with ode <%d>", url, clientResp.StatusCode)
		response <- res
		return
	}
	res.Resp = clientResp
	res.Err = nil
	response <- res
}
