package NCMB

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Query struct {
	ncmb      *NCMB
	className string
	where     map[string]interface{}
	limit     int
	skip      int
	count     bool
	order     string
	include   string
}

func (query *Query) EqualTo(key string, value interface{}) *Query {
	if query.where == nil {
		query.where = make(map[string]interface{})
	}
	query.where[key] = value
	return query
}

func (query *Query) NotEqualTo(key string, value interface{}) *Query {
	return query.setOperand(key, "$ne", value)
}

func (query *Query) LessThan(key string, value interface{}) *Query {
	return query.setOperand(key, "$lt", value)
}

func (query *Query) LessThanOrEqualTo(key string, value interface{}) *Query {
	return query.setOperand(key, "$lte", value)
}

func (query *Query) GreaterThan(key string, value interface{}) *Query {
	return query.setOperand(key, "$gt", value)
}

func (query *Query) GreaterThanOrEqualTo(key string, value interface{}) *Query {
	return query.setOperand(key, "$gte", value)
}

func (query *Query) In(key string, value interface{}) *Query {
	return query.setOperand(key, "$in", value)
}

func (query *Query) NotIn(key string, value interface{}) *Query {
	return query.setOperand(key, "$nin", value)
}

func (query *Query) Exists(key string, value interface{}) *Query {
	return query.setOperand(key, "$exists", value)
}

func (query *Query) RegularExpression(key string, value string) *Query {
	return query.setOperand(key, "$regex", value)
}

func (query *Query) InArray(key string, value interface{}) *Query {
	return query.setOperand(key, "$inArray", value)
}

func (query *Query) NotInArray(key string, value interface{}) *Query {
	return query.setOperand(key, "$ninArray", value)
}

func (query *Query) AllInArray(key string, value interface{}) *Query {
	return query.setOperand(key, "$all", value)
}

func (query *Query) Near(key string, value GeoPoint) *Query {
	return query.setOperand(key, "$nearSphere", value)
}

func (query *Query) WithinKilometers(key string, value GeoPoint, distance float64) *Query {
	return query.setOperand(key, "$nearSphere", value).setOperand(key, "$maxDistanceInKilometers", distance)
}

func (query *Query) WithinMiles(key string, value GeoPoint, distance float64) *Query {
	return query.setOperand(key, "$nearSphere", value).setOperand(key, "$maxDistanceInMiles", distance)
}

func (query *Query) WithinRadians(key string, value GeoPoint, distance float64) *Query {
	return query.setOperand(key, "$nearSphere", value).setOperand(key, "$maxDistanceInRadians", distance)
}

func (query *Query) WithinSquare(key string, southWest GeoPoint, northEast GeoPoint) *Query {
	box := []GeoPoint{southWest, northEast}
	value := map[string]interface{}{"$box": box}
	return query.setOperand(key, "$within", value)
}

type RelatedItem interface {
	GetClassName() string
	GetObjectId() string
}

func (query *Query) RelatedTo(object RelatedItem, key string) *Query {
	if query.where == nil {
		query.where = make(map[string]interface{})
	}
	query.where["$relatedTo"] = map[string]interface{}{
		"object": map[string]interface{}{"__type": "Pointer", "className": object.GetClassName(), "objectId": object.GetObjectId()},
		"key":    key,
	}
	return query
}

func (query *Query) Limit(num int) *Query {
	query.limit = num
	return query
}

func (query *Query) Include(key string) *Query {
	query.include = key
	return query
}

func (query *Query) Order(key string, descending ...bool) *Query {
	if len(descending) > 0 && descending[0] {
		key = "-" + key
	} else {
		query.order = key
	}
	return query
}

func (query *Query) setOperand(key string, operand string, value interface{}) *Query {
	if query.where == nil {
		query.where = make(map[string]interface{})
	}
	var operandMap map[string]interface{}
	if query.where[key] == nil {
		operandMap = make(map[string]interface{})
	} else {
		operandMap = query.where[key].(map[string]interface{})
	}
	operandMap[operand] = changeValue(value)
	query.where[key] = operandMap
	return query
}

func changeValue(value interface{}) interface{} {
	if reflect.TypeOf(value).Name() == "Time" {
		val := value.(time.Time).UTC()
		return map[string]interface{}{"__type": "Date", "iso": val.Format("2006-01-02T15:04:05.999Z0700")}
	}
	if reflect.TypeOf(value).Name() == "Item" {
		return map[string]interface{}{
			"__type":    "Pointer",
			"className": value.(*Item).ClassName,
			"objectId":  value.(*Item).ObjectId,
		}
	}
	return value
}

func (query *Query) FetchAll() ([]Item, error) {
	queries := make(map[string]interface{})
	if len(query.where) > 0 {
		queries["where"] = query.where
	}
	if query.limit > 0 {
		if query.limit > 1000 {
			return nil, fmt.Errorf("limit is over 1000")
		}
		queries["limit"] = query.limit
	}
	if query.skip > 0 {
		queries["skip"] = query.skip
	}
	if query.count {
		queries["count"] = 1
	}
	if query.order != "" {
		queries["order"] = query.order
	}
	if query.include != "" {
		queries["include"] = query.include
	}
	fmt.Println(queries)
	request := Request{ncmb: query.ncmb}
	params := ExecOptions{}
	params.ClassName = query.className
	params.Queries = &queries
	data, err := request.Gets(params)
	if err != nil {
		return nil, err
	}
	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	fmt.Println(results)
	ary, err := json.Marshal(results["results"])
	if err != nil {
		return nil, err
	}
	var aryResults []map[string]interface{}
	err = json.Unmarshal(ary, &aryResults)
	if err != nil {
		return nil, err
	}
	var items []Item
	for _, value := range aryResults {
		item := Item{ncmb: query.ncmb, ClassName: query.className}
		item.Sets(value)
		items = append(items, item)
	}
	return items, nil
}
