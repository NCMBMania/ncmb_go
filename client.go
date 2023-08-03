package NCMB

// import "fmt"

type NCMB struct {
	applicationKey string
	clientKey      string
	sessionToken   string
}

func Initialize(applicationKey string, clientKey string) NCMB {
	ncmb := NCMB{applicationKey: applicationKey, clientKey: clientKey}
	return ncmb
}

func (ncmb *NCMB) DataStore(className string) DataStore {
	dataStore := DataStore{ncmb: ncmb, className: className}
	return dataStore
}

func (ncmb *NCMB) Query(className string) Query {
	return Query{ncmb: ncmb, className: className}
}
