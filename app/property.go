package app

type PropertyProvider struct {
	properties map[string]interface{}
}

func (pp *PropertyProvider) GetProperty(property string) (value interface{}, exists bool) {
	value, exists = pp.properties[property]
	return value, exists
}

func (pp *PropertyProvider) SetProperty(property string, value interface{}) {
	pp.properties[property] = value
}