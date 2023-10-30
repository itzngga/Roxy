package compress

import (
	"bytes"
	"compress/flate"
	"github.com/andybalholm/brotli"
	"github.com/goccy/go-json"
	"io"
)

func MarshallBrotli(v interface{}) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var fBuf bytes.Buffer
	fw, err := flate.NewWriter(&fBuf, 9)
	if err != nil {
		return nil, err
	}

	_, err = fw.Write(result)
	if err != nil {
		return nil, err
	}

	err = fw.Close()
	if err != nil {
		return nil, err
	}

	var bBuf bytes.Buffer
	bw := brotli.NewWriterLevel(&bBuf, brotli.BestCompression)

	_, err = bw.Write(fBuf.Bytes())
	if err != nil {
		return nil, err
	}

	err = bw.Close()
	if err != nil {
		return nil, err
	}

	defer func() {
		result = nil
		fBuf.Reset()
		bBuf.Reset()
	}()

	return bBuf.Bytes(), nil
}

func UnmarshallBrotli[T any](data []byte, v T) error {
	var dataReader = bytes.NewReader(data)
	br := brotli.NewReader(dataReader)

	var flateDecoded bytes.Buffer
	_, err := io.Copy(&flateDecoded, br)
	if err != nil {
		return err
	}

	dataReader = bytes.NewReader(flateDecoded.Bytes())
	flateReader := flate.NewReader(dataReader)

	var decoded bytes.Buffer
	_, err = io.Copy(&decoded, flateReader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decoded.Bytes(), v)
	if err != nil {
		return err
	}

	defer func() {
		br = nil
		dataReader = nil
		flateReader.Close()
		decoded.Reset()
		flateDecoded.Reset()
		data = nil
	}()

	return nil
}
