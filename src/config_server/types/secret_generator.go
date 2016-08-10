package types

import (
    "crypto/rand"
    "math/big"
)

type secretGenerator struct {
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func NewSecretGenerator() ValueGenerator {
    return secretGenerator{}
}

func (secretGenerator) Generate(parameters interface{}) (interface{}, error) {
    max := big.NewInt(int64(len(letterRunes)))
    runes := make([]rune, 20)

    for i := range runes {
        idx, err := rand.Int(rand.Reader, max)
        if err != nil {
            return nil, err
        }

        runes[i] = letterRunes[idx.Int64()]
    }

    return string(runes), nil
}
