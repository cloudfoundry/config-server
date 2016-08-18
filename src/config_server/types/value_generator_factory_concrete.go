package types

import (
    "errors"
    "fmt"
    "config_server/config"
)

type valueGeneratorConcrete struct {
    config config.ServerConfig
}

func NewValueGeneratorConcrete(config config.ServerConfig) valueGeneratorConcrete {
    return valueGeneratorConcrete{config: config}
}

func (vgc valueGeneratorConcrete) GetGenerator(valueType string) (ValueGenerator, error) {
    switch valueType {
    case "password":
        return NewPasswordGenerator(), nil
    case "certificate":
        x509Loader := NewX509Loader()
        return NewCertificateGenerator(vgc.config, x509Loader), nil
    default:
        return nil, errors.New(fmt.Sprintf("Unsupported value type: %s", valueType))
    }
}
