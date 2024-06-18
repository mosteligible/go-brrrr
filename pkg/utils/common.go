package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
		fmt.Println("Cannot parse response body to struct!", err.Error())
		return false
	}
	if _, ok := parsedResp["error"]; ok {
		log.Printf("Error response: %v", parsedResp)
		return false
	}
	return true
}

func ParseDateStringtoTime(datestr string) (time.Time, error) {
	var parsed time.Time
	parsed, err := time.Parse(time.DateOnly, datestr)
	if err == nil {
		return parsed, nil
	}
	log.Printf("Parsing with date only <%s> failed, err: %s", time.DateOnly, err.Error())
	parsed, err = time.Parse(time.DateTime, datestr)
	if err == nil {
		return parsed, nil
	}
	log.Printf("Parsing with date only <%s> failed, err: %s", time.DateTime, err.Error())
	return GetZero[time.Time](), err
}

func GetKeysOfMap[T any](kv map[string]T) []string {
	retval := []string{}
	for k := range kv {
		retval = append(retval, k)
	}
	return retval
}

func WriteToJsonFile[T any](content T, filename string, indent bool) error {
	var err error
	var jsonData []byte
	if indent {
		jsonData, err = json.MarshalIndent(content, "", "  ")
	} else {
		jsonData, err = json.Marshal(content)
	}

	if err != nil {
		return err
	}
	os.WriteFile("./out/"+time.Now().Format(time.DateTime)+filename, jsonData, os.ModePerm)

	return nil
}
