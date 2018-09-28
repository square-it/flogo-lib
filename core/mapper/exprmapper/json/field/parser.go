package field

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

type MappingField struct {
	hasSpecialField bool
	hasArray        bool
	fields          []string
	s               *scanner.Scanner
}

func NewMappingField(hasSpecialField, hasArray bool, fields []string) *MappingField {
	return &MappingField{hasSpecialField: hasSpecialField, hasArray: hasArray, fields: fields}
}

func ParseMappingField(mRef string) (*MappingField, error) {
	g := &MappingField{}
	//Remove any . at beginning
	if strings.HasPrefix(mRef, ".") {
		mRef = mRef[1:]
	}
	fmt.Println(mRef)
	err := g.Start(mRef)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (m *MappingField) Getfields() []string {
	return m.fields
}

func (m *MappingField) HasSepcialField() bool {
	return m.hasSpecialField
}

func (m *MappingField) HasArray() bool {
	return m.hasArray
}

func (m *MappingField) paserName() error {
	fieldName := ""
	switch ch := m.s.Scan(); ch {
	case '.':
		return m.Parser()
	case '[':
		//Done
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
		m.s.Mode = scanner.ScanInts
		nextAfterBracket := m.s.Scan()
		if nextAfterBracket == '"' || nextAfterBracket == '\'' {
			//Special charactors
			m.s.Mode = scanner.ScanIdents
			return m.handleSpecialField()
		} else {
			//HandleArray
			m.hasArray = true
			m.fields[len(m.fields)-1] = m.fields[len(m.fields)-1] + "[" + m.s.TokenText() + "]"
			//m.handleArray()
			ch := m.s.Scan()
			if ch != ']' {
				return fmt.Errorf("Inliad array format")
			}
			m.s.Mode = scanner.ScanIdents
			return m.Parser()
		}
	case scanner.EOF:
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
	default:
		fieldName = fieldName + m.s.TokenText()
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
		return m.Parser()
	}

	return nil
}

func (m *MappingField) handleSpecialField() error {
	m.hasSpecialField = true
	fieldName := ""
	run := true

	for run {
		switch ch := m.s.Scan(); ch {
		case '"', '\'':
			nextAfterQuotes := m.s.Scan()
			if nextAfterQuotes == ']' {
				//end specialfield, startover
				m.fields = append(m.fields, fieldName)
				run = false
				return m.Parser()
			} else {
				fieldName = fieldName + m.s.TokenText()
			}
		default:
			fieldName = fieldName + m.s.TokenText()
		}
	}
	return nil
}

func (m *MappingField) Parser() error {
	switch ch := m.s.Scan(); ch {
	case '.':
		return m.paserName()
	case '[':
		m.s.Mode = scanner.ScanInts
		nextAfterBracket := m.s.Scan()
		if nextAfterBracket == '"' || nextAfterBracket == '\'' {
			//Special charactors
			m.s.Mode = scanner.ScanIdents
			return m.handleSpecialField()
		} else {
			//HandleArray
			m.fields[len(m.fields)-1] = m.fields[len(m.fields)-1] + "[" + m.s.TokenText() + "]"
			//m.handleArray()
			ch := m.s.Scan()
			if ch != ']' {
				return fmt.Errorf("Inliad array format")
			}
			m.hasArray = true
			m.s.Mode = scanner.ScanIdents
			return m.Parser()
		}
	case scanner.EOF:
		//Done
		return nil
	default:
		m.fields = append(m.fields, m.s.TokenText())
		return m.paserName()
	}
	return nil
}

func (m *MappingField) Start(jsonPath string) error {
	m.s = new(scanner.Scanner)
	m.s.IsIdentRune = IsIdentRune
	m.s.Init(strings.NewReader(jsonPath))
	m.s.Mode = scanner.ScanIdents
	return m.Parser()
}

func IsIdentRune(ch rune, i int) bool {

	//except ', ", [, ] when all should be ok.
	//if ch == '\'' || ch == '[' || ch == ']' || ch == '.' {
	//	return false
	//}
	return ch == '$' || ch == '-' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0

}
