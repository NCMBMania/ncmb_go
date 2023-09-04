package NCMB

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestItemSave(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
}

func TestItemDelete(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	bol, err = hello.Delete()
	if err != nil {
		t.Errorf("hello.Delete() = %T, %s", bol, err)
	}
}

func TestItemFetch(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	hello := ncmb.Item("Hello")
	msg := "Hello, World!"
	hello.Set("msg", msg)
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	message, err := hello.GetString("msg")
	if err != nil {
		t.Errorf("hello.GetString(msg) = %s, %s", message, err)
	}
	if message != msg {
		t.Errorf("hello.GetString(msg) = %s, want %s", message, msg)
	}
	hello2 := ncmb.Item("Hello")
	hello2.ObjectId = hello.ObjectId
	bol, err = hello2.Fetch()
	if err != nil {
		t.Errorf("hello2.Fetch() = %T, %s", bol, err)
	}
	str1, err := hello.GetString("msg")
	str2, err := hello2.GetString("msg")
	if str1 != str2 {
		t.Errorf("hello2.GetString(msg) = %s, want %s", str2, str1)
	}
	if hello.ObjectId != hello2.ObjectId {
		t.Errorf("hello.ObjectId = %s, want %s", hello.ObjectId, hello2.ObjectId)
	}
	bol, err = hello.Delete()
	if err != nil {
		t.Errorf("hello.Delete() = %T, %s", bol, err)
	}
	bol, err = hello2.Fetch()
	if err == nil {
		t.Errorf("hello2 deleted failed.")
	}
	if err.Error() != "NCMBError: E404001, No data available." {
		t.Errorf("hello2.Fetch() = %T, %s", bol, err)
	}
}

func TestItemUpdate(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	msg := "Hello, World! 2"
	hello.Set("msg", msg).Set("num", 100)
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	hello.Fetch()
	str, err := hello.GetString("msg")
	if err != nil {
		t.Errorf("hello.GetString(msg) = %s, %s", str, err)
	}
	if str != msg {
		t.Errorf("hello.GetString(msg) = %s, want %s", str, msg)
	}
	hello.Delete()
}

func TestItemSaveData(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	hello := ncmb.Item("Hello")
	hello.Set("string", "Hello, World!")
	hello.Set("number", 139.77421902)
	hello.Set("boolean", true)
	hello.Set("array", []interface{}{"a", "b", "c"})
	hello.Set("object", map[string]string{"key": "value", "key2": "value2"})
	hello.Set("null", nil)
	hello.Set("date", time.Now())
	geo := ncmb.GeoPoint(35.698683, 139.77421902)
	hello.Set("geo", geo)
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if bol != true {
		t.Errorf("hello.Save() = %T, want true", bol)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello.Fetch()
	str, err := hello.GetString("string")
	if err != nil {
		t.Errorf("hello.GetString(string) = %s, %s", str, err)
	}
	if str != "Hello, World!" {
		t.Errorf("hello.GetString(string) = %s, want Hello, World!", str)
	}
	num, err := hello.GetNumber("number")
	if err != nil {
		t.Errorf("hello.GetNumber(number) = %f, %s", num, err)
	}
	if num != 139.77421902 {
		t.Errorf("hello.GetNumber(number) = %f, want 139.77421902", num)
	}
	bol, err = hello.GetBool("boolean")
	if err != nil {
		t.Errorf("hello.GetBool(boolean) = %T, %s", bol, err)
	}
	if bol != true {
		t.Errorf("hello.GetBool(boolean) = %T, want true", bol)
	}
	arr, err := hello.GetArray("array")
	if err != nil {
		t.Errorf("hello.GetArray(array) = %T, %s", arr, err)
	}
	date, err := hello.GetDate("date")
	if err != nil {
		t.Errorf("hello.GetDate(date) = %T, %s", date, err)
	}
	if reflect.TypeOf(date).Name() != "Time" {
		t.Errorf("hello.GetDate(date) = %T, want time.Time", date)
	}
	geo, err = hello.GetGeoPoint("geo")
	if err != nil {
		t.Errorf("hello.GetGeoPoint(geo) = %T, %s", geo, err)
	}
	if geo.Latitude != 35.698683 {
		t.Errorf("hello.GetGeoPoint(geo) = %f, want 35.698683", geo.Latitude)
	}
	hello.Delete()
}

func TestItemInclude(t *testing.T) {
	ncmb := SetUpNCMB()
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello2 := ncmb.Item("Hello2")
	hello2.Set("msg", "Hello, World!!!")
	bol, err = hello2.Save()
	if err != nil {
		t.Errorf("hello2.Save() = %T, %s", bol, err)
	}
	if hello2.ObjectId == "" {
		t.Errorf("hello2.ObjectId = %s, want not empty", hello2.ObjectId)
	}
	hello.Set("hello2", hello2)
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	query := ncmb.Query("Hello")
	query.Include("hello2")
	query.EqualTo("objectId", hello.ObjectId)
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	item := items[0]
	h, err := item.GetItem("hello2")
	if err != nil {
		t.Errorf("item.GetItem(hello2) = %T, %s", h, err)
	}
	if h.ObjectId != hello2.ObjectId {
		t.Errorf("Expected: %s, Actual: %s", hello2.ObjectId, h.ObjectId)
	}
}
func TestItemIncrement(t *testing.T) {
	ncmb := SetUpNCMB()
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	hello.Set("num", 3)
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello.Increment("num", 2)
	hello.Increment("num", 2)
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	hello.Fetch()
	num, err := hello.GetNumber("num")
	if err != nil {
		t.Errorf("hello.GetNumber(num) = %T, %s", num, err)
	}
	if num != 7 {
		t.Errorf("hello.GetNumber(num) = %f, want 7", num)
	}
}

func TestItemAdd(t *testing.T) {
	ncmb := SetUpNCMB()
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	hello.Set("array", []int{1})
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello.Add("array", 2)
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	hello.Fetch()
	ary := hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 := ary.([]interface{})
	if len(ary2) != 2 {
		t.Errorf("len(ary2) = %d, want 2", len(ary2))
	}
	if int(ary2[0].(float64)) != 1 {
		t.Errorf("ary2[0] = %d, want 1", ary2[0])
	}
	hello.Add("array", 3)
	hello.Add("array", 4)
	hello.Save()
	hello.Fetch()
	ary = hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 = ary.([]interface{})
	if len(ary2) != 4 {
		t.Errorf("len(ary2) = %d, want 4", len(ary2))
	}
	if int(ary2[3].(float64)) != 4 {
		t.Errorf("ary2[3] = %d, want 4", ary2[3])
	}
	hello.Delete()
}

func TestItemAddUnique(t *testing.T) {
	ncmb := SetUpNCMB()
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	hello.Set("array", []interface{}{1})
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello.AddUnique("array", "2")
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	hello.Fetch()
	ary := hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 := ary.([]interface{})
	if len(ary2) != 2 {
		t.Errorf("len(ary2) = %d, want 2", len(ary2))
	}
	if int(ary2[0].(float64)) != 1 {
		t.Errorf("ary2[0] = %d, want 1", ary2[0])
	}
	hello.AddUnique("array", 3)
	hello.AddUnique("array", 4)
	hello.AddUnique("array", 3)
	hello.AddUnique("array", "2")
	hello.Save()
	hello.Fetch()
	ary = hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 = ary.([]interface{})
	if len(ary2) != 4 {
		t.Errorf("len(ary2) = %d, want 4", len(ary2))
	}
	if int(ary2[3].(float64)) != 4 {
		t.Errorf("ary2[3] = %d, want 4", ary2[3])
	}
	hello.Delete()
}

func TestItemRemove(t *testing.T) {
	ncmb := SetUpNCMB()
	hello := ncmb.Item("Hello")
	hello.Set("msg", "Hello, World!")
	hello.Set("array", []interface{}{1, 3})
	bol, err := hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	if hello.ObjectId == "" {
		t.Errorf("hello.ObjectId = %s, want not empty", hello.ObjectId)
	}
	hello.Add("array", 2)
	bol, err = hello.Save()
	if err != nil {
		t.Errorf("hello.Save() = %T, %s", bol, err)
	}
	hello.Fetch()
	ary := hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 := ary.([]interface{})
	if len(ary2) != 3 {
		t.Errorf("len(ary2) = %d, want 3", len(ary2))
	}
	if int(ary2[0].(float64)) != 1 {
		t.Errorf("ary2[0] = %d, want 1", ary2[0])
	}
	hello.Remove("array", 3)
	hello.Remove("array", 2)
	hello.Save()
	hello.Fetch()
	ary = hello.Get("array")
	if reflect.TypeOf(ary).Kind() != reflect.Slice {
		t.Errorf("hello.Get(array) = %f, want array", ary)
	}
	ary2 = ary.([]interface{})
	if len(ary2) != 1 {
		t.Errorf("len(ary2) = %d, want 1, %s", len(ary2), ary2)
	}
	if int(ary2[0].(float64)) != 1 {
		t.Errorf("ary2[1] = %d, want 1", ary2)
	}
	hello.Delete()
}
