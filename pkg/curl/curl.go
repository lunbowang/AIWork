package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func PostRequest(tokenStr, url string, requestBody any) ([]byte, error) {
	return sendRequest(tokenStr, url, "POST", requestBody)
}

func DeleteRequest(tokenStr, url string, requestBody any) ([]byte, error) {
	return sendRequest(tokenStr, url, "DELETE", nil)
}

func PutRequest(tokenStr, url string, requestBody any) ([]byte, error) {
	return sendRequest(tokenStr, url, "PUT", requestBody)
}

func GetRequest(tokenStr, urls string, queryParams map[string]any) ([]byte, error) {
	// 拼接查询参数
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, fmt.Sprintf("%v", value))
		}
		urls = urls + "?" + values.Encode()
	}
	return sendRequest(tokenStr, urls, "GET", nil)
}

func sendRequest(tokenStr, url, method string, requestBody any) ([]byte, error) {

	var (
		body []byte
		err  error
	)

	if requestBody != nil {
		body, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if len(tokenStr) > 0 {
		req.Header.Set("Authorization", tokenStr)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBuffer := new(bytes.Buffer)
	_, err = responseBuffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBuffer.Bytes(), nil
}
