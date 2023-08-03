package NCMB

import (
	"fmt"
	"reflect"
	"github.com/harakeishi/gats"
	"time"
	"encoding/json"
)

type Item struct {
	dataStore *DataStore
	ObjectId string
	fields map[string]interface{}
}

func (item *Item) Set(key string, value interface{}) *Item {
	if item.fields == nil {
		item.fields = make(map[string]interface{})
	}
	switch key {
	case "objectId":
		val, err := gats.ToString(value)
		if err == nil {
			item.ObjectId = val
		}
	case "createDate", "updateDate":
		val, err := gats.ToString(value)
		if err != nil {
			break
		}
		date, err := time.Parse("2006-01-02T15:04:05.999Z0700", val)
		if err == nil {
			item.fields[key] = date
		}
	default:
		item.fields[key] = value
	}
	return item
}

func (item *Item) GetString(key string, defaultValue ...string) (string, error) {
	value := item.fields[key]
	if (value == nil) {
		if (defaultValue != nil && len(defaultValue) > 0) {
			return defaultValue[0], nil
		}
		return "", fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf("") {
		return "", fmt.Errorf("%s is not string (%s)", key, reflect.TypeOf(value))
	}
	return value.(string), nil
}

func (item *Item) GetInt(key string, defaultValue ...int) (int, error) {
	value := item.fields[key]
	if (value == nil) {
		if (defaultValue != nil && len(defaultValue) > 0) {
			return defaultValue[0], nil
		}
		return -1, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf(0) {
		return -1, fmt.Errorf("%s is not int (%s)", key, reflect.TypeOf(value))
	}
	return value.(int), nil
}

func (item *Item) GetDate(key string, defaultValue ...time.Time) (time.Time, error) {
	value := item.fields[key]
	if (value == nil) {
		if (defaultValue != nil && len(defaultValue) > 0) {
			return defaultValue[0], nil
		}
		return time.Now(), fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf(time.Now()) {
		return time.Now(), fmt.Errorf("%s is not time.Time (%s)", key, reflect.TypeOf(value))
	}
	return value.(time.Time), nil
}

func (item *Item) GetFloat(key string, defaultValue ...float64) (float64, error) {
	value := item.fields[key]
	if (value == nil) {
		if (defaultValue != nil && len(defaultValue) > 0) {
			return defaultValue[0], nil
		}
		return 0.001, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf(0.001) {
		return 0.001, fmt.Errorf("%s is not float64 (%s)", key, reflect.TypeOf(value))
	}
	return value.(float64), nil
}

func (item *Item) Save() (bool, error) {
	if item.ObjectId == "" {
		return item.Create()
	} else {
		return item.Update()
	}
}

func (item *Item) Sets(hash map[string]interface{}) *Item {
	for key, value := range hash {
		item.Set(key, value)
	}
	return item
}

func (item *Item) Create() (bool, error) {
	request := Request{ncmb: item.dataStore.ncmb}
	data, err := request.Post(item.dataStore.className, item.fields)
	if err != nil {
		return false, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(data, &hash)
	if err != nil {
		return false, err
	}
	item.Sets(hash)
	return true, nil
}

func (item *Item) Update() (bool, error) {
	// TODO
	return true, nil
}
