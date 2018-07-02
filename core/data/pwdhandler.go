package data

var passwordDecoder PasswordValueDecoder

type PasswordValueDecoder interface {
	DecodeValue(value interface{}) (string, error)
}


func SetPasswordValueDecoder(pwdResolver PasswordValueDecoder ) {
	passwordDecoder = pwdResolver
}

func GetPasswordValueDecoder() PasswordValueDecoder {
	if passwordDecoder == nil {
		passwordDecoder = &defaultPasswordValueDecoder{}
	}
	return passwordDecoder
}

type defaultPasswordValueDecoder struct {

}

func (defaultResolver *defaultPasswordValueDecoder) DecodeValue(value interface{}) (string, error) {
	if value != nil {
		return value.(string), nil
	}
	return "", nil
}



