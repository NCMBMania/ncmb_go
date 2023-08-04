# NCMB SDK for Golang

## Usage

### Install

```
go get -u github.com/NCMBMania/ncmb_go
```

### Import

```go
import (
    "github.com/NCMBMania/ncmb_go"
)
```

### Initialize

```go
ncmb := NCMB.Initialize("YOUR_APPLICATION_KEY", "YOUR_CLIENT_KEY")
```

### DataStore

#### Save item

```go
item := ncmb.Item("Hello")
item.Set("msg1", "Hello").Set("msg2", "World")
item.Set("num", 100)      // Int
item.Set("float", 1.23)   // Float

// Save item
bol, err := item.Save()
```

#### Get item data

```go
msg1, err := item.GetString("msg1"))
msg3, err := item.GetString("msg3", "default") // 2nd argument is default value
bol, err := item.GetBool("bool") // true or false
num, err := item.GetNumber("num") // float64
ary, err := item.GetArray("ary") // []interface{}
obj, err := item.GetMap("map") // interface{}{}
geo, err := item.GetGeoPoint("geo") // GeoPoint
date, err := item.GetDate("date") // time.Time
other, err := item.Get("null") // interface{}
item.ObjectId // string
```

### Query

#### Fetch data

```go
query := ncmb.Query("Hello")
query.EqualTo("msg1", "Hello")
items, err := query.FetchAll()
if err != nil {
	fmt.Println(err)
}
fmt.Println(items[0].GetDate("createDate"))
```

### GeoPoint

#### Add GeoPoint as data

```go
item := ncmb.Item("Hello")
geo := ncmb.GeoPiint(35.6585805, 139.7454329)
item.Set("geo", geo)
bol, err := item.Save()
```

## License

MIT.

