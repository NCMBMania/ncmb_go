package NCMB

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type Request struct {
	ncmb *NCMB
}

type ExecOptions struct {
	ClassName         string
	Fields            *map[string]interface{}
	ObjectId          *string
	Queries           *map[string]interface{}
	AdditionalHeaders *map[string]string
	Path              *string
	Multipart         bool
	IsScript          bool
	Bytes             *[]byte
}

type NCMBError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (request *Request) Post(options ...ExecOptions) ([]byte, error) {
	return request.Exec("POST", options[0])
}

func (request *Request) Put(options ...ExecOptions) ([]byte, error) {
	return request.Exec("PUT", options[0])
}

func (request *Request) Gets(options ...ExecOptions) ([]byte, error) {
	return request.Exec("GET", options[0])
}

func (request *Request) Get(options ...ExecOptions) ([]byte, error) {
	return request.Exec("GET", options[0])
}

func (request *Request) Delete(options ...ExecOptions) ([]byte, error) {
	return request.Exec("DELETE", options[0])
}

func (request *Request) Data(data *map[string]interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	if (*data)["acl"] != nil {
		acl, err := (*data)["acl"].(Acl).ToJSON()
		if err != nil {
			return nil, err
		}
		(*data)["acl"] = acl
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (request *Request) Exec(method string, params ExecOptions) ([]byte, error) {
	opts := UrlOptions{ObjectId: params.ObjectId, DefinePath: params.Path, Queries: params.Queries}
	s := Signature{ncmb: request.ncmb, IsScript: params.IsScript}
	s.Initialize()
	sig, err := s.Generate(method, params.ClassName, opts)
	if err != nil {
		return nil, err
	}
	url, err := s.Url(params.ClassName, opts)
	if err != nil {
		return nil, err
	}
	headers := s.Headers(sig)
	if !params.Multipart {
		headers["Content-Type"] = "application/json"
	}
	if request.ncmb.SessionToken != "" {
		headers["X-NCMB-Apps-Session-Token"] = request.ncmb.SessionToken
	}
	if params.AdditionalHeaders != nil {
		for key, value := range *params.AdditionalHeaders {
			headers[key] = value
		}
	}
	client := &http.Client{}
	data := &bytes.Buffer{}
	if method == "POST" || method == "PUT" {
		if params.Multipart {
			mw := multipart.NewWriter(data)
			part := make(textproto.MIMEHeader)
			mimeType := http.DetectContentType(*params.Bytes)
			mimeType = strings.Split(mimeType, ";")[0]
			fmt.Println(mimeType)
			part.Set("Content-Type", mimeType)
			part.Set("Content-Disposition", `form-data; name="file"; filename="file"`)
			pw, err := mw.CreatePart(part)
			if err != nil {
				return nil, err
			}
			io.Copy(pw, bytes.NewReader(*params.Bytes))
			if err != nil {
				return nil, err
			}
			acl := (*params.Fields)["acl"]
			if acl != nil {
				value, err := acl.(Acl).ToJSON()
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(value)
				if err != nil {
					return nil, err
				}
				err = mw.WriteField("acl", string(data))
				if err != nil {
					return nil, err
				}
			}
			headers["Content-Type"] = mw.FormDataContentType()
			err = mw.Close()
			if err != nil {
				return nil, err
			}
		} else {
			d, err := request.Data(params.Fields)
			if err != nil {
				return nil, err
			}
			data = bytes.NewBuffer(d)
		}
		fmt.Println(params.Fields)
	}
	fmt.Println(url)
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	fmt.Println(headers)
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
