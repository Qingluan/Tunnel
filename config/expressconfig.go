package config

type ExpressConfig struct {
	Raw      interface{}
	Multi    bool
	Protocol string
	First    []byte
}

func DefaultConfig() *ExpressConfig {
	return &ExpressConfig{
		Multi:    false,
		Protocol: "tls",
	}
}

func Config(s interface{}) ExpressConfig {
	return ExpressConfig{
		Raw: s,
	}
}

func (exp ExpressConfig) Check(setconfig *ExpressConfig) {
	if exp.Raw != nil {
		switch exp.Raw.(type) {
		case string:
			exp.Protocol = exp.Raw.(string)
			setconfig.Protocol = exp.Raw.(string)
			exp.Raw = nil
		case bool:
			exp.Multi = exp.Raw.(bool)

			setconfig.Multi = exp.Raw.(bool)
			exp.Raw = nil
		case []byte:
			exp.First = exp.Raw.([]byte)
			setconfig.First = exp.Raw.([]byte)
			exp.Raw = nil
		}
	}
}

func ParseConfigs(configs ...interface{}) (parsedConfig *ExpressConfig) {
	defaultconfig := DefaultConfig()
	for _, c := range configs {
		Config(c).Check(defaultconfig)
	}
	return defaultconfig
}
