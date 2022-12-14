package gofast

import (
	"bytes"
	"github.com/goccy/go-json"
	"log"
	"mime/multipart"
	"net/textproto"

	"github.com/valyala/fasthttp"
)

type RequestEncoder func(req *fasthttp.Request, in interface{}) error

type ResponseDecoder func(resp *fasthttp.Response, out interface{}) error

var JSONEncoder = func(req *fasthttp.Request, in interface{}) error {
	req.Header.SetContentType("application/json")
	return json.NewEncoder(req.BodyWriter()).Encode(in)
}

var JSONDecoder = func(resp *fasthttp.Response, out interface{}) error {
	if err := json.Unmarshal(resp.Body(), out); err != nil {
		log.Printf("[gofast] response decode failed - code: %v, body: %v", resp.StatusCode(), string(resp.Body()))
		return err
	}
	return nil
}

var URLEncoder = func(req *fasthttp.Request, in interface{}) error {
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	for k, v := range in.(Body) {
		args.Set(k, v)
	}
	if _, err := args.WriteTo(req.BodyWriter()); err != nil {
		return err
	}
	req.Header.SetContentType("application/x-www-form-urlencoded")
	return nil
}

var BodyEncoder = func(req *fasthttp.Request, in interface{}) error {
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})

	res, _ := json.Marshal(in)
	part.Write(res)

	writer.Close()

	req.MultipartForm()

	req.Header.SetContentType("multipart/form-data")
	return nil
}

var TextDecoder = func(resp *fasthttp.Response, out interface{}) error {
	s := out.(*string)
	*s = string(resp.Body())
	return nil
}

var ByteDecoder = func(resp *fasthttp.Response, out interface{}) error {
	b := out.(*[]byte)
	*b = resp.Body()
	return nil
}
