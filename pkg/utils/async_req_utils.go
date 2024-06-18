package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mosteligible/go-brrrr/pkg/types"
)

var UnsuccessCount = 0
var SuccessCount = 0

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
	// log.Printf("Sending request to %s\nheaders: %v", url, headers)
	client := &http.Client{Timeout: 30 * time.Second}
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
		if err != nil {
			fmt.Println("error making post request")
			os.Exit(2)
		}
	default:
		res.Err = fmt.Errorf("Method not supported: %s", method)
		response <- res
		return
	}
	setHeaders(req, &headers)
	clientResp, err = client.Do(req)
	if err != nil {
		// log.Printf("Error making request to %s - clientResp: %v", url, clientResp)
		UnsuccessCount++
		res.Err = err
		response <- res
		return
	}
	if clientResp.StatusCode > 399 {
		log.Printf("Client <%s> Responded with code <%d>", url, clientResp.StatusCode)
		response <- res
		return
	}
	res.Resp = clientResp
	res.Err = nil
	response <- res
	SuccessCount++
}
