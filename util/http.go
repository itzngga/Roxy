package util

import (
	"io"
	"net/http"
)

func DoHTTPRequest(method string, url string) ([]byte, error) {
	var client = &http.Client{}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
