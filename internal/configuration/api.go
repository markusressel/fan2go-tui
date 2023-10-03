package configuration

type ApiConfig struct {
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
}
