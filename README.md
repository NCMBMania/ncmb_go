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

### User

#### Sign up

```go
user, err := ncmb.SignUpByAccount(userName, password)
```

#### Log in by username

```go
user, err := ncmb.Login(userName, password)
```

#### Log in by email address

```go
user, err := ncmb.LoginWithMailAddress(mailAddress, password)
```

#### Request registration email

```go
bol, err := ncmb.RequestSignUpEmail(mailAddress)
```

#### Reset password (only email auth)

```go
bol, err := ncmb.RequestPasswordReset(mailAddress)
```

#### Logout

```go
bol, err := ncmb.Logout()
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

Supported operands.

- func EqualTo(key string, value interface{})
- func NotEqualTo(key string, value interface{})
- func LessThan(key string, value interface{})
- func LessThanOrEqualTo(key string, value interface{})
- func GreaterThan(key string, value interface{})
- func GreaterThanOrEqualTo(key string, value interface{})
- func In(key string, value interface{})
- func NotIn(key string, value interface{})
- func Exists(key string, value interface{})
- func RegularExpression(key string, value string)
- func InArray(key string, value interface{})
- func NotInArray(key string, value interface{})
- func AllInArray(key string, value interface{})
- func Near(key string, value GeoPoint)
- func WithinKilometers(key string, value GeoPoint, distance float64)
- func WithinMiles(key string, value GeoPoint, distance float64)
- func WithinRadians(key string, value GeoPoint, distance float64)
- func WithinSquare(key string, southWest GeoPoint, northEast GeoPoint)

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

