package NCMB

type DataStore struct {
	ncmb *NCMB
	className string
}

func (dataStore *DataStore) Item() Item {
	item := Item{dataStore: dataStore, fields: make(map[string]interface{})}
	return item
}
