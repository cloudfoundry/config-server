package types

import (
	"crypto/rand"
	"math/big"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type passwordGenerator struct {
}

type passwordParams struct {
	AllowedCharacters string `yaml:"allowed_characters"`
	Length            int    `yaml:"length"`
}

var defaultLetterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

const DefaultPasswordLength = 20

func NewPasswordGenerator() ValueGenerator {
	return passwordGenerator{}
}

func (passwordGenerator) Generate(parameters interface{}) (interface{}, error) {
	letterRunes := defaultLetterRunes
	length := DefaultPasswordLength
	if parameters != nil {
		var params passwordParams
		err := objToStruct(parameters, &params)
		if err != nil {
			return nil, errors.Error("Failed to generate password, parameters are invalid.")
		}
		if params.AllowedCharacters != "" {
			letterRunes = []rune(params.AllowedCharacters)
		}
		if params.Length > 10 {
			length = params.Length
		}
	}
	lengthLetterRunes := big.NewInt(int64(len(letterRunes)))
	passwordRunes := make([]rune, length)
	for i := range passwordRunes {
		index, err := rand.Int(rand.Reader, lengthLetterRunes)
		if err != nil {
			return nil, err
		}

		passwordRunes[i] = letterRunes[index.Int64()]
	}

	return string(passwordRunes), nil
}
