package types

import (
	"fmt"
	"net/http"
)

type RespHolder struct {
	Resp *http.Response
	Err  error
}

type Parameters map[string]interface{}

func (p Parameters) String() string {
	retval := "Parameters{\n"

	for key, value := range p {
		row := fmt.Sprintf("    %s: %v\n", key, value)
		retval = fmt.Sprintf("%s%s", retval, row)
	}

	return fmt.Sprintf("%s}\n", retval)
}
