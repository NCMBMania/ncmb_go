package NCMB

// import "fmt"

type NCMB struct {
	ApplicationKey string
	ClientKey      string
	SessionToken   string
}

func Initialize(applicationKey string, clientKey string) NCMB {
	ncmb := NCMB{ApplicationKey: applicationKey, ClientKey: clientKey}
	return ncmb
}

func (ncmb *NCMB) Item(className string) Item {
	return Item{ncmb: ncmb, className: className}
}

func (ncmb *NCMB) Query(className string) Query {
	return Query{ncmb: ncmb, className: className}
}

func (ncmb *NCMB) GeoPoint(latitude float64, longitude float64) GeoPoint {
	return GeoPoint{Type: "GeoPoint", Latitude: latitude, Longitude: longitude}
}
