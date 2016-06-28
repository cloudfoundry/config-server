package store

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

type EncoderDecoder struct {
}

type val struct {
	Value interface{}
}

func NewEncoderDecoder() EncoderDecoder {
	return EncoderDecoder{}
}

func (EncoderDecoder) Encode(toEncode interface{}) (string, error) {
	toEncode = val{toEncode}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(toEncode)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func (EncoderDecoder) Decode(encodedStr string) (interface{}, error) {
	var ret val

	b64DecodedBytes, _ := base64.StdEncoding.DecodeString(encodedStr)
	buff := bytes.Buffer{}
	buff.Write(b64DecodedBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret.Value, nil
}
