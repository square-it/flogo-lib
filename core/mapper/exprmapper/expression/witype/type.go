package witype

type TYPE int

const (
	STRING TYPE = iota
	INT64
	REF
	ARRAYREF
	FLOAT
	FUNCTION
	EXPRESSION
	BOOL
)

func (t TYPE) String() string {
	switch t {
	case STRING:
		return "string"
	case INT64:
		return "int64"
	case REF:
		return "ref"
	case ARRAYREF:
		return "arrayRef"
	case FLOAT:
		return "float"
	case FUNCTION:
		return "function"
	case EXPRESSION:
		return "expression"
	}
	return ""
}
