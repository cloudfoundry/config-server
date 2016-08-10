package types

import (
    "errors"
    "fmt"
)

type valueGeneratorConcrete struct {
}

func NewValueGeneratorConcrete() valueGeneratorConcrete {
    return valueGeneratorConcrete{}
}

func (valueGeneratorConcrete) GetGenerator(valueType string) (ValueGenerator, error) {
    switch valueType {
    case "password":
        return NewPasswordGenerator(), nil
    default:
        return nil, errors.New(fmt.Sprintf("Unsupported value type: %s", valueType))
    }
}
