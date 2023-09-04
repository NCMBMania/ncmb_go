package NCMB

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestDeleteData(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	query := ncmb.Query("Query")
	query.Limit(1000)
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	for _, item := range items {
		bol, err := item.Delete()
		if err != nil {
			t.Errorf("item.Delete() = %T, %s", bol, err)
		}
	}
}

func TestQueryFetchAll(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	query := ncmb.Query("Query")
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) > 0 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}
}

func TestQueryFetchAll2(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("godotenv.Load() = %s", err)
		return
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	// create 5 items
	for i := 0; i < 5; i++ {
		item := ncmb.Item("Query")
		item.Set("msg", "Hello, World!")
		item.Set("index", i)
		item.Save()
	}
	query := ncmb.Query("Query")
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 5 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}
	// Delete items
	for _, item := range items {
		bol, err := item.Delete()
		if err != nil {
			t.Errorf("item.Delete() = %T, %s", bol, err)
		}
	}
}

func SetUpNCMB() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func SetUp() NCMB {
	ncmb := SetUpNCMB()
	// create 5 items
	t := time.Now()
	for i := 0; i < 5; i++ {
		item := ncmb.Item("Query")
		item.Set("msg", "Hello, World!")
		item.Set("index", i)
		item.Set("bol", i%2 == 0)
		item.Set("date", t.AddDate(0, 0, i))
		item.Save()
	}
	return ncmb
}

func TearDown(ncmb NCMB) {
	query := ncmb.Query("Query")
	query.Limit(1000)
	items, err := query.FetchAll()
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		bol, err := item.Delete()
		if err != nil {
			panic(err)
		}
		if bol != true {
			panic("item.Delete() = false")
		}
	}
}

func TestQueryFetchNumber(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	query.GreaterThan("index", 2)
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 2 {
		t.Errorf("len(items) = %d, want 2", len(items))
	}
	TearDown(ncmb)
}

func TestQueryFetchNumber2(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	query.GreaterThan("index", 1)
	query.LessThan("index", 4)
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 2 {
		t.Errorf("len(items) = %d, want 2", len(items))
	}
	TearDown(ncmb)
}

func TestQueryFetchNumber3(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	query.GreaterThanOrEqualTo("index", 1)
	query.LessThanOrEqualTo("index", 4)
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 4 {
		t.Errorf("len(items) = %d, want 4", len(items))
	}
	TearDown(ncmb)
}

func TestQueryFetchBol(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	query.EqualTo("bol", true)
	items1, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items1, err)
	}
	query = ncmb.Query("Query")
	query.Limit(1000)
	query.NotEqualTo("bol", true)
	items2, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items2, err)
	}
	if len(items1)+len(items2) != 5 {
		t.Errorf("len(items) = %d, want 5", len(items1)+len(items2))
	}
	TearDown(ncmb)
}

func TestQueryFetchDate(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	now := time.Now()
	query.GreaterThan("date", now.Add(time.Hour*1))
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 4 {
		t.Errorf("len(items) = %d, want 4", len(items))
	}
	TearDown(ncmb)
}

func TestQueryFetchDate2(t *testing.T) {
	ncmb := SetUp()
	query := ncmb.Query("Query")
	query.Limit(1000)
	now := time.Now().Add(time.Hour * 1)
	query.GreaterThan("date", now)
	query.LessThan("date", now.AddDate(0, 0, 2))
	items, err := query.FetchAll()
	if err != nil {
		t.Errorf("query.FetchAll() = %T, %s", items, err)
	}
	if len(items) != 2 {
		t.Errorf("len(items) = %d, want 2", len(items))
	}
	TearDown(ncmb)
}

func TestQueryRelatedTo(t *testing.T) {
	ncmb := SetUpNCMB()
	item := ncmb.Item("Query")
	item.ObjectId = "aaa"
	query := ncmb.Query("Query")
	query.RelatedTo(&item, "hoge")
	TearDown(ncmb)
}
