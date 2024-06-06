package types

import (
	"net/http"
)

type RespHolder struct {
	Resp *http.Response
	Err  error
}
