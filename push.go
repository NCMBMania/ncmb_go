package NCMB

import (
	"fmt"
	"reflect"
)

type Push struct {
	ncmb *NCMB
	Item
}

func (push *Push) Save() (bool, error) {
	_, err := push.Item.GetDate("deliveryTime")
	if err != nil {
		immediateDeliveryFlag, err := push.Item.GetBool("immediateDeliveryFlag")
		if err != nil || immediateDeliveryFlag == false {
			return false, fmt.Errorf("immediateDeliveryFlag or deliveryTime is required")
		}
	}
	targets := push.Item.Get("target")
	if targets == nil {
		return false, fmt.Errorf("target is required. %s", targets)
	}
	if reflect.TypeOf(targets).Kind() != reflect.Slice {
		return false, fmt.Errorf("target is not array. %s", targets)
	}
	if len(targets.([]string)) == 0 {
		return false, fmt.Errorf("target is required. %s", targets)
	}
	return push.Item.Save()
}
