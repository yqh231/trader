package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type httpClient struct {
	sync.RWMutex

	client  *http.Client
	headers map[string]string
	root    string
}

type httpResp struct {
	resp *http.Response
}

func newHttpClient(root string, timeOut time.Duration) *httpClient {
	return &httpClient{
		client:  &http.Client{Timeout: timeOut * time.Second},
		headers: make(map[string]string),
		root:    root,
	}
}

func (client *httpClient) setHeaders(headers map[string]string) {
	client.Lock()
	defer client.Unlock()

	for k, v := range headers {
		client.headers[k] = v
	}
}

func (client *httpClient) Get(url string, parameters map[string]string) (*httpResp, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	req, err = http.NewRequest("GET", client.root+url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	for k, v := range parameters {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	resp, err = client.client.Do(req)
	if err != nil {
		return nil, err
	}

	return &httpResp{resp: resp}, nil
}

func (client *httpClient) Post(url string, parameters map[string]interface{}) (*httpResp, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
		body []byte
	)

	body, err = json.Marshal(parameters)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest("POST", client.root+url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err = client.client.Do(req)
	if err != nil {
		return nil, err
	}

	return &httpResp{resp: resp}, nil
}

func (r *httpResp) toString() (string, error) {
	response, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (r *httpResp) unmarshal(v interface{}) error {
	response, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(response, v)
	if err != nil {
		return err
	}

	return nil
}
