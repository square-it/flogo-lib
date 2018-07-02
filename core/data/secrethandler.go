package data

var secretDecoder SecretValueDecoder

type SecretValueDecoder interface {
	DecodeValue(value interface{}) (string, error)
}


func SetSecretValueDecoder(pwdResolver SecretValueDecoder ) {
	secretDecoder = pwdResolver
}

func GetSecretValueDecoder() SecretValueDecoder {
	if secretDecoder == nil {
		secretDecoder = &defaultSecretValueDecoder{}
	}
	return secretDecoder
}

type defaultSecretValueDecoder struct {

}

func (defaultResolver *defaultSecretValueDecoder) DecodeValue(value interface{}) (string, error) {
	if value != nil {
		return value.(string), nil
	}
	return "", nil
}



