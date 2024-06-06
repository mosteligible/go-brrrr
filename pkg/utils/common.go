package utils

import (
	"encoding/json"
	"io"
	"log"
)

func GetZero[T any]() T {
	var retval T
	return retval
}

func ParseResponseBodyToStruct[T any](response *io.ReadCloser) (T, error) {
	var resp T
	decoder := json.NewDecoder(*response)
	if err := decoder.Decode(&resp); err != nil {
		return GetZero[T](), err
	}
	return resp, nil
}

func CheckErrorInResponse(responseBody *io.ReadCloser) bool {
	var parsedResp map[string]interface{}
	parsedResp, err := ParseResponseBodyToStruct[map[string]interface{}](responseBody)
	if err != nil {
		return false
	}
	if _, ok := parsedResp["error"]; ok {
		log.Printf("Error response: %v", parsedResp)
		return false
	}
	return true
}
