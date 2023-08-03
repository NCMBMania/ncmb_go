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

func (request *Request) Post(class_name string, fields map[string]interface{}, multipart ...bool) (map[string]interface{}, error) {
	params := ExecOptions{Multipart: multipart != nil, Fields: &fields}
	return request.Exec("POST", class_name, params)
}

func (request *Request) Put() (map[string]interface{}, error) {
	// TODO
	return nil, nil
}

func (request *Request) Get() (map[string]interface{}, error) {
	// TODO
	return nil, nil
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

func (request *Request) Exec(method string, className string, params ExecOptions) (map[string]interface{}, error) {
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
	if method == "POST" {
		data, err := request.Data(params.Fields)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
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
		// fmt.Println(res.Body)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var response map[string]interface{}
		json.Unmarshal(body, &response)
		return response, nil
	} else {
		res := map[string]interface{}{}
		return res, nil
	}
}