package pab

// Config represent pab configuration
type Config struct {
	Name        string
	Description string
	Adapter     string
	Storage     string
	HTTPAddr    string
	Env         map[string]interface{}
}

// NewConfig create pab configuration with default values
func NewConfig() *Config {
	return &Config{
		Name:        "pab",
		Description: "a chatter bot",
		Adapter:     "shell",
		Storage:     "memory",
		HTTPAddr:    ":8001",
		Env:         make(map[string]interface{}),
	}
}
