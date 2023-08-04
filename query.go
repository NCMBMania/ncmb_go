package NCMB

import (
	"fmt"
	"encoding/json"
)

type Query struct {
	ncmb *NCMB
	className string
	Where map[string]interface{}
	Limit int
	Skip int
	Count bool
	Order string
	Include string
}

func (query *Query) EqualTo(key string, value interface{}) *Query {
	if query.Where == nil {
		query.Where = make(map[string]interface{})
	}
	query.Where[key] = value
	return query
}

func (query *Query) FetchAll() ([]Item, error) {
	queries := make(map[string]interface{})
	queries["where"] = query.Where
	if query.Limit > 0 {
		if query.Limit > 1000 {
			return nil, fmt.Errorf("limit is over 1000")
		}
		queries["limit"] = query.Limit
	}
	if query.Skip > 0 {
		queries["skip"] = query.Skip
	}
	if query.Count {
		queries["count"] = 1
	}
	request := Request{ncmb: query.ncmb}
	data, err := request.Gets(query.className, queries)
	if err != nil {
		return nil, err
	}
	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
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
		item := Item{ncmb: query.ncmb, className: query.className}
		item.Sets(value)
		items = append(items, item)
	}
	return items, nil
}

