package NCMB

import (
	"bytes"
	"net/http"
	"encoding/json"
	// "fmt"
	"io/ioutil"
)

type Request struct {
	ncmb *NCMB
}

type ExecOptions struct {
	Fields *map[string]interface{}
	ObjectId *string
	Queries *map[string]interface{}
	AdditionalHeaders *map[string]string
	Path *string
	Multipart bool
	IsScript bool
}

func (request *Request) Post(class_name string, fields map[string]interface{}, multipart ...bool) ([]byte, error) {
	params := ExecOptions{Multipart: multipart != nil, Fields: &fields}
	return request.Exec("POST", class_name, params)
}

func (request *Request) Put() (map[string]interface{}, error) {
	// TODO
	return nil, nil
}

func (request *Request) Get(class_name string, queries map[string]interface{}, multipart ...bool) ([]byte, error) {
	params := ExecOptions{Multipart: multipart != nil, Queries: &queries}
	return request.Exec("GET", class_name, params)
}

func (request *Request) Delete() (bool, error) {
	// TODO
	return true, nil
}

func (request *Request) Data(data *map[string]interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	/*
	for key, value := range *data {
		fmt.Println(key, value)
	}
	*/
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (request *Request) Exec(method string, className string, params ExecOptions) ([]byte, error) {
	opts := UrlOptions{ObjectId: params.ObjectId, DefinePath: params.Path, Queries: params.Queries}
	s := Signature{ncmb: request.ncmb, IsScript: params.IsScript}
	s.Initialize()
	sig, err := s.Generate(method, className, opts)
	if err != nil {
		return nil, err
	}
	url, err := s.Url(className, opts)
	if err != nil {
		return nil, err
	}
	headers := s.Headers(sig)
	if params.AdditionalHeaders != nil {
		for key, value := range *params.AdditionalHeaders {
			headers[key] = value
		}
	}
	client := &http.Client{}
	data := new(bytes.Buffer)
	if method == "POST" || method == "PUT" {
		d, err := request.Data(params.Fields)
		if err != nil {
			return nil, err
		}
		data = bytes.NewBuffer(d)
	}

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
