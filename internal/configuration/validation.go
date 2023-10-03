package configuration

func Validate(configPath string) error {
	return validateConfig(&CurrentConfig, configPath)
}

func validateConfig(config *Configuration, path string) error {
	return nil
}
