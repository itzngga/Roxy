package gofast

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type ErrorHandler func(resp *fasthttp.Response) error

var defaultErrorHandler = func(resp *fasthttp.Response) error {
	return fmt.Errorf("code: %v, body: %v", resp.StatusCode(), string(resp.Body()))
}
