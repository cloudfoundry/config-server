package types

import (
	"config_server/config"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type ValueGeneratorConcrete struct {
	config config.ServerConfig
}

func NewValueGeneratorConcrete(config config.ServerConfig) ValueGeneratorConcrete {
	return ValueGeneratorConcrete{config: config}
}

func (vgc ValueGeneratorConcrete) GetGenerator(valueType string) (ValueGenerator, error) {
	switch valueType {
	case "password":
		return NewPasswordGenerator(), nil
	case "certificate":
		x509Loader := NewX509Loader()
		return NewCertificateGenerator(vgc.config, x509Loader), nil
	default:
		return nil, errors.Errorf("Unsupported value type: %s", valueType)
	}
}
