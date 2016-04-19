package data

type Type int

const (
	STRING Type = iota
	INTEGER
	NUMBER
	BOOLEAN
	OBJECT
	ARRAY
	PARAMS
)

var types = [...]string {
	"string",
	"integer",
	"number",
	"boolean",
	"object",
	"array",
	"params",
}

var typeMap = map[string]Type{
	"string": STRING,
	"integer": INTEGER,
	"number": NUMBER,
	"boolean": BOOLEAN,
	"object": OBJECT,
	"array": ARRAY,
	"params": PARAMS,
}

func (t Type) String() string {
	return types[t]
}

func ToType(typeStr string) (Type, bool) {
	dataType, found :=  typeMap[typeStr]

	return dataType, found
}
