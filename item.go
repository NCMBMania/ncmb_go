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
	ClassName string
	ObjectId  string
	fields    map[string]interface{}
}

type ItemDate struct {
	Type string `json:"__type"`
	Iso  string `json:"iso"`
}

type AddRemoveOperation struct {
	Op      string        `json:"__op"`
	Objects []interface{} `json:"objects"`
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
		fmt.Println(val)
		item.fields[key] = ItemDate{Type: "Date", Iso: val}
	case "acl":
		if value == nil {
			break
		}
		if reflect.TypeOf(value).Name() == "Acl" {
			item.fields[key] = value.(Acl)
		} else {
			acl := Acl{}
			for key, value := range value.(map[string]interface{}) {
				for key2, value2 := range value.(map[string]interface{}) {
					acl.setAccess(key, key2, value2.(bool))
				}
			}
			item.SetAcl(acl)
		}
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
	if reflect.TypeOf(value).Kind() == reflect.Map {
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
	if reflect.TypeOf(value).Kind() != reflect.Slice {
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

func (item *Item) GetItem(key string, defaultValue ...Item) (Item, error) {
	value := item.fields[key]
	if value == nil {
		if defaultValue != nil && len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return Item{}, fmt.Errorf("key %s is not found", key)
	}
	if reflect.TypeOf(value).Name() == "Item" {
		return value.(Item), nil
	}
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return Item{}, fmt.Errorf("%s is not Map (%s)", key, reflect.TypeOf(value))
	}
	valueMap := value.(map[string]interface{})
	if valueMap["__type"] != "Object" {
		return Item{}, fmt.Errorf("%s is not Item format (%s)", key, valueMap)
	}
	className := valueMap["className"].(string)
	delete(valueMap, "__type")
	delete(valueMap, "className")
	i := item.ncmb.Item(className)
	i.Sets(valueMap)
	return i, nil
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
	options := ExecOptions{}
	options.ClassName = item.ClassName
	fields := item.Fields()
	options.Fields = &fields
	data, err := request.Post(options)
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
	options := ExecOptions{}
	options.ClassName = item.ClassName
	options.ObjectId = &item.ObjectId
	fields := item.Fields()
	options.Fields = &fields
	data, err := request.Put(options)
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

func (item *Item) SetAcl(acl Acl) *Item {
	return item.Set("acl", acl)
}

func (item *Item) GetAcl() (Acl, error) {
	params := item.fields["acl"]
	if params == nil {
		acl := Acl{}
		acl.SetPublicReadAccess(true).SetPublicWriteAccess(true)
		return acl, nil
	}
	if reflect.TypeOf(params).Name() == "Acl" {
		return params.(Acl), nil
	}
	if reflect.TypeOf(params).Kind() != reflect.Map {
		return Acl{}, fmt.Errorf("acl is not Map (%s)", reflect.TypeOf(params))
	}
	valueMap := params.(map[string]map[string]bool)
	acl := Acl{}
	acl.Sets(valueMap)
	return acl, nil
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
			val := value.(time.Time).UTC()
			hash[key] = ItemDate{Type: "Date", Iso: val.Format("2006-01-02T15:04:05.000Z")}
		} else if reflect.TypeOf(value).Name() == "Item" {
			val := value.(Item)
			hash[key] = val.ToPointer()
		} else {
			hash[key] = value
		}
	}
	return hash
}

func (item *Item) Delete() (bool, error) {
	request := Request{ncmb: item.ncmb}
	params := ExecOptions{}
	params.ClassName = item.ClassName
	params.ObjectId = &item.ObjectId
	data, err := request.Delete(params)
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
	params := ExecOptions{}
	params.ClassName = item.ClassName
	params.ObjectId = &item.ObjectId
	data, err := request.Get(params)
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

func (item *Item) ToPointer() map[string]interface{} {
	return map[string]interface{}{
		"__type":    "Pointer",
		"className": item.ClassName,
		"objectId":  item.ObjectId,
	}
}

func (item *Item) GetClassName() string {
	return item.ClassName
}

func (item *Item) GetObjectId() string {
	return item.ObjectId
}

func (item *Item) Increment(key string, amount ...int) *Item {
	amountInt := 1
	if amount != nil && len(amount) > 0 {
		amountInt = amount[0]
	}
	value := item.Get(key)
	if value != nil && reflect.TypeOf(value).Kind() == reflect.Map {
		num := value.(map[string]interface{})["amount"].(int)
		amountInt += num
		fmt.Println(num, amountInt)
	}
	item.Set(key, map[string]interface{}{
		"__op":   "Increment",
		"amount": amountInt,
	})
	return item
}

func (item *Item) Add(key string, object interface{}) *Item {
	if object == nil {
		return item
	}
	value := item.Get(key)
	if item.ObjectId != "" {
		if value != nil && reflect.TypeOf(value).Kind() == reflect.Map {
			values := value.(AddRemoveOperation)
			if values.Op == "Add" {
				values.Objects = append(values.Objects, object)
				item.Set(key, values)
				return item
			}
		}
		val := AddRemoveOperation{}
		val.Op = "Add"
		val.Objects = []interface{}{object}
		item.Set(key, val)
	} else {
		if value != nil && reflect.TypeOf(value).Kind() != reflect.Map {
			value = append(value.([]interface{}), object)
			return item.Set(key, value)
		}
		item.Set(key, []interface{}{object})
		return item
	}
	return item
}

func (item *Item) AddUnique(key string, object interface{}) *Item {
	if object == nil {
		return item
	}
	value := item.Get(key)
	values := AddRemoveOperation{}
	if value == nil || (reflect.TypeOf(value).Kind() == reflect.Map && value.(map[string]interface{})["__op"] != "AddUnique") {
		values = AddRemoveOperation{
			Op:      "AddUnique",
			Objects: []interface{}{},
		}
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		values = AddRemoveOperation{
			Op:      "AddUnique",
			Objects: []interface{}{},
		}
		for _, v := range value.([]interface{}) {
			values.Objects = append(values.Objects, v)
		}
	} else {
		values = value.(AddRemoveOperation)
	}
	if slices.Index(values.Objects, object) == -1 {
		values.Objects = append(values.Objects, object)
	}
	item.Set(key, values)
	return item
}

func (item *Item) Remove(key string, object interface{}) *Item {
	if object == nil {
		return item
	}
	value := item.Get(key)
	values := AddRemoveOperation{}
	if value == nil || (reflect.TypeOf(value).Kind() == reflect.Map && value.(map[string]interface{})["__op"] != "Remove") {
		values = AddRemoveOperation{
			Op:      "Remove",
			Objects: []interface{}{},
		}
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		values = AddRemoveOperation{
			Op:      "Remove",
			Objects: []interface{}{},
		}
	} else {
		values = value.(AddRemoveOperation)
	}
	values.Objects = append(values.Objects, object)
	item.Set(key, values)
	return item
}
