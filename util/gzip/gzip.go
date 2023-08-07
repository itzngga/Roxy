package gzip

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/klauspost/compress/gzip"
	"io"
)

func MarshallGzip(v interface{}) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	var zw = gzip.NewWriter(&buf)

	_, err = zw.Write(result)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	defer func() {
		result = nil
		buf.Reset()
	}()

	return buf.Bytes(), nil
}

func UnmarshallGzip[T any](data []byte, v T) error {
	var dataReader = bytes.NewReader(data)
	zr, err := gzip.NewReader(dataReader)
	if err != nil {
		return err
	}

	var decoded bytes.Buffer
	_, err = io.Copy(&decoded, zr)
	if err != nil {
		return err
	}

	if err := zr.Close(); err != nil {
		return err
	}

	err = json.Unmarshal(decoded.Bytes(), v)
	if err != nil {
		return err
	}

	defer func() {
		dataReader = nil
		decoded.Reset()
		data = nil
	}()

	return nil
}
