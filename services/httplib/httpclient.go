package httplib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// An interface to send http request
type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
	Post(url string, body io.Reader) (resp *http.Response, err error)
}

// Holds a session with current user's token
type Session struct {
	Token        string
	HttpClient   HttpClient
	HttpHeaders  map[string]string
	BaseEndPoint string
}

// Request an API with `path` the GET method, and optional `params`
func (session *Session) Get(path string, params Params) (response []byte, err error) {
	urlStr := session.getUrl(path, params)
	response, err = session.sendGetRequest(urlStr)
	return

	//res, err = MakeResult(response)
	//return

}

// Request an API with `path` and `data` the POST method
func (session *Session) Post(path string, urlParams, data Params) (response []byte, err error) {
	urlStr := session.getUrl(path, urlParams)

	response, err = session.sendPostRequest(urlStr, data)

	return

	//res, err = MakeResult(response)
	//return
}

// Get generic url with `path` and optional `params`
func (session *Session) getUrl(path string, params Params) string {
	buf := &bytes.Buffer{}
	buf.WriteString(path)

	if params != nil {
		buf.WriteRune('?')

		params.EncodeUrlParam(buf)
	}

	return buf.String()
}

func (session *Session) sendGetRequest(url string) ([]byte, error) {
	var request *http.Request
	var err error

	request, err = http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	for k, v := range session.HttpHeaders {
		request.Header.Set(k, v)
	}
	return session.sendRequest(request)
}

func (session *Session) sendPostRequest(url string, params Params) ([]byte, error) {
	buf := &bytes.Buffer{}

	params.EncodeRequestBody(buf)
	var request *http.Request
	var err error

	request, err = http.NewRequest("POST", url, buf)

	if err != nil {
		return nil, err
	}

	for k, v := range session.HttpHeaders {
		request.Header.Set(k, v)
	}

	request.Header.Set("Content-Type", "application/json")

	return session.sendRequest(request)
}

func (session *Session) sendRequest(request *http.Request) ([]byte, error) {
	var response *http.Response
	var err error

	if session.HttpClient == nil {
		response, err = http.DefaultClient.Do(request)
	} else {
		response, err = session.HttpClient.Do(request)
	}

	if err != nil {
		return nil, fmt.Errorf("账号中心服务异常. %v", err)
	}

	defer response.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, response.Body)

	if err != nil {
		return nil, fmt.Errorf("账号中心服务异常. %v", err)
	}

	if response.StatusCode != 200 {
		ret := make(map[string]interface{})
		if err := json.Unmarshal(buf.Bytes(), &ret); err == nil {
			return nil, fmt.Errorf("%v", ret["error"])
		} else {
			return nil, err
		}

	}

	return buf.Bytes(), nil
}
