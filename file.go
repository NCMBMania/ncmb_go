package NCMB

import (
	"encoding/json"
	"fmt"
	"time"
)

type File struct {
	ncmb  *NCMB
	Bytes []byte
	Item
}

func (file *File) GetDate(key string, defaultValue ...time.Time) (time.Time, error) {
	return file.Item.GetDate(key, defaultValue...)
}

func (file *File) Update() (bool, error) {
	params := ExecOptions{}
	params.ClassName = "files"
	fileName, err := file.Item.GetString("fileName")
	if err != nil {
		return false, err
	}
	params.ObjectId = &fileName
	acl, err := file.Item.GetAcl()
	if err != nil {
		return false, err
	}
	params.Fields = &map[string]interface{}{"acl": acl}
	params.Multipart = false
	request := Request{ncmb: file.ncmb}
	result, err := request.Put(params)
	if err != nil {
		return false, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(result, &hash)
	if err != nil {
		return false, err
	}
	file.Sets(hash)
	file.Item.Sets(hash)
	return true, nil
}

func (file *File) Upload() (bool, error) {
	params := ExecOptions{}
	params.ClassName = "files"
	fileName, err := file.Item.GetString("fileName")
	if err != nil {
		return false, err
	}
	params.ObjectId = &fileName
	acl, err := file.Item.GetAcl()
	if err != nil {
		return false, err
	}
	params.Fields = &map[string]interface{}{"acl": acl}
	params.Bytes = &file.Bytes
	params.Multipart = true
	request := Request{ncmb: file.ncmb}
	result, err := request.Post(params)
	if err != nil {
		return false, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(result, &hash)
	if err != nil {
		return false, err
	}
	file.Sets(hash)
	file.Item.Sets(hash)
	return true, nil
}

func (file *File) Download() ([]byte, error) {
	request := Request{ncmb: file.ncmb}
	params := ExecOptions{}
	params.ClassName = "files"
	fileName, err := file.Item.GetString("fileName")
	if err != nil {
		return nil, err
	}
	params.ObjectId = &fileName
	data, err := request.Get(params)
	if err != nil {
		return nil, err
	}
	file.Bytes = data
	return data, nil
}

func (file *File) Delete() (bool, error) {
	request := Request{ncmb: file.ncmb}
	params := ExecOptions{}
	params.ClassName = "files"
	fileName, err := file.Item.GetString("fileName")
	if err != nil {
		return false, err
	}
	params.ObjectId = &fileName
	data, err := request.Delete(params)
	if err != nil {
		return false, err
	}
	if string(data) == "" {
		return true, nil
	}
	return false, fmt.Errorf("delete error, %s", data)
}
