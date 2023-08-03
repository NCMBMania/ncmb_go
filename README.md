# NCMB SDK for Golang

## Usage

### Initialize

```go
ncmb := NCMB.Initialize("YOUR_APPLICATION_KEY", "YOUR_CLIENT_KEY")
```

### DataStore

#### Get DataStore

```go
datastore := ncmb.DataStore("YOUR_CLASS_NAME")
```

#### Save item

```go
item := datastore.Item()
item.Set("msg1", "Hello").Set("msg2", "World")
item.Set("num", 100)      // Int
item.Set("float", 1.23)   // Float

// Save item
item.Save()
```

#### Get item data

```go
msg1, err := item.GetString("msg1"))
msg1, err := item.GetString("msg3", "default") // 2nd argument is default value
item.GetInt("num")
item.GetFloat("float")
item.ObjectId
```

## License

MIT.

