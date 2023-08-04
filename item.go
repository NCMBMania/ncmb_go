package NCMB

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/harakeishi/gats"
	"golang.org/x/exp/slices"
)

type Item struct {
	ncmb      *NCMB
	className string
	ObjectId  string
	fields    map[string]interface{}
}

type ItemDate struct {
	Type string `json:"__type"`
	Iso  string `json:"iso"`
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
		item.fields[key] = ItemDate{Type: "Date", Iso: val}
	default:
		item.fields[key] = value
	}
	return item
}

func (item *Item) Get(key string) interface{} {
	return item.fields[key]
}

func (item *Item) GetString(key string, defaultValue ...string) (string, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf("") {
		return "", fmt.Errorf("%s is not string (%s)", key, reflect.TypeOf(value))
	}
	return value.(string), nil
}

func (item *Item) GetDate(key string, defaultValue ...time.Time) (time.Time, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return time.Now(), fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf(time.Now()) {
		// Object?
		if reflect.TypeOf(value).Kind() != reflect.Map {
			return time.Now(), fmt.Errorf("%s is not time.Time (%s)", key, reflect.TypeOf(value))
		}
		val := value.(map[string]interface{})
		if val["__type"] != "Date" {
			return time.Now(), fmt.Errorf("%s is not Date format (%s)", key, val)
		}
		date, err := time.Parse("2006-01-02T15:04:05.999Z0700", val["iso"].(string))
		if err != nil {
			return time.Now(), err
		}
		return date, nil
	}
	return value.(time.Time), nil
}

func (item *Item) GetArray(key string, defaultValue ...[]interface{}) ([]interface{}, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return []interface{}{}, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf([]interface{}{}) {
		return []interface{}{}, fmt.Errorf("%s is not []interface{} (%s)", key, reflect.TypeOf(value))
	}
	return value.([]interface{}), nil
}

func (item *Item) GetMap(key string, defaultValue ...map[string]interface{}) (map[string]interface{}, error) {
	value := item.fields[key]
	nullValue := make(map[string]interface{})
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return nullValue, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return nullValue, fmt.Errorf("%s is not interface{}{} (%s)", key, reflect.TypeOf(value))
	}
	return value.(map[string]interface{}), nil
}

func (item *Item) GetBool(key string, defaultValue ...bool) (bool, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return false, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value) != reflect.TypeOf(false) {
		return false, fmt.Errorf("%s is not bool (%s)", key, reflect.TypeOf(value))
	}
	return value.(bool), nil
}

func (item *Item) GetGeoPoint(key string, defaultValue ...GeoPoint) (GeoPoint, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return GeoPoint{}, fmt.Errorf("key is not found")
	}
	if reflect.TypeOf(value).Name() == "GeoPoint" {
		return value.(GeoPoint), nil
	}
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return GeoPoint{}, fmt.Errorf("%s is not Map (%s)", key, reflect.TypeOf(value))
	}
	valueMap := value.(map[string]interface{})
	if valueMap["__type"] != "GeoPoint" {
		return GeoPoint{}, fmt.Errorf("%s is not GeoPoint format (%s)", key, valueMap)
	}
	latitude, longitude := valueMap["latitude"].(float64), valueMap["longitude"].(float64)
	return GeoPoint{Latitude: latitude, Longitude: longitude}, nil
}

func (item *Item) GetNumber(key string, defaultValue ...float64) (float64, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
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
	request := Request{ncmb: item.ncmb}
	data, err := request.Post(item.className, item.Fields())
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
	request := Request{ncmb: item.ncmb}
	data, err := request.Put(item.className, item.ObjectId, item.Fields())
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

func (item *Item) Fields() map[string]interface{} {
	if item.fields == nil {
		return make(map[string]interface{})
	}
	hash := make(map[string]interface{})
	for key, value := range item.fields {
		if slices.Index([]string{"objectId", "createDate", "updateDate"}, key) > -1 {
			continue
		}
		if value == nil {
			hash[key] = nil
			continue
		}
		if reflect.TypeOf(value).Name() == "Time" {
			hash[key] = ItemDate{Type: "Date", Iso: value.(time.Time).Format("2006-01-02T15:04:05.000Z")}
		} else {
			hash[key] = value
		}
	}
	return hash
}

func (item *Item) Delete() (bool, error) {
	request := Request{ncmb: item.ncmb}
	data, err := request.Delete(item.className, item.ObjectId)
	if err != nil {
		return false, err
	}
	if string(data) != "" {
		return false, fmt.Errorf("delete error %s", string(data))
	}
	return true, nil
}

func (item *Item) Fetch() (bool, error) {
	request := Request{ncmb: item.ncmb}
	data, err := request.Get(item.className, item.ObjectId)
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
