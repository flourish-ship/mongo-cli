package httplib

import (
	"io"
	"net/http"
)

type BaseService struct {
	Token string
}

func (this *BaseService) DoReq(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(header) > 0 {
		req.Header = header
	}

	req.Header.Set("Authorization", this.Token)

	client := &http.Client{}
	return client.Do(req)
}
