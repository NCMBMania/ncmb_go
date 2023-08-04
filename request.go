package NCMB

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	ncmb *NCMB
}

type ExecOptions struct {
	Fields            *map[string]interface{}
	ObjectId          *string
	Queries           *map[string]interface{}
	AdditionalHeaders *map[string]string
	Path              *string
	Multipart         bool
	IsScript          bool
}

type NCMBError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (request *Request) Post(className string, fields map[string]interface{}, multipart ...bool) ([]byte, error) {
	params := ExecOptions{Multipart: multipart != nil, Fields: &fields}
	return request.Exec("POST", className, params)
}

func (request *Request) Put(className string, objectId string, fields map[string]interface{}, multipart ...bool) ([]byte, error) {
	params := ExecOptions{Multipart: multipart != nil, ObjectId: &objectId, Fields: &fields}
	return request.Exec("PUT", className, params)
}

func (request *Request) Gets(className string, queries map[string]interface{}, multipart ...bool) ([]byte, error) {
	params := ExecOptions{Multipart: multipart != nil, Queries: &queries}
	return request.Exec("GET", className, params)
}

func (request *Request) Get(className string, objectId string) ([]byte, error) {
	params := ExecOptions{ObjectId: &objectId}
	return request.Exec("GET", className, params)
}

func (request *Request) Delete(className string, objectId string) ([]byte, error) {
	params := ExecOptions{ObjectId: &objectId}
	return request.Exec("DELETE", className, params)
}

func (request *Request) Data(data *map[string]interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
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
	if !params.Multipart {
		headers["Content-Type"] = "application/json"
	}
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
	fmt.Println(data)
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
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_, err = json.Marshal(string(body))
	if err == nil {
		var ncmbError NCMBError
		err = json.Unmarshal(body, &ncmbError)
		if err == nil && ncmbError.Code != "" {
			return nil, fmt.Errorf("NCMBError: %s, %s", ncmbError.Code, ncmbError.Error)
		}
	}
	return body, nil
}
