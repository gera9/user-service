package config

type Config struct {
	App struct {
		Name    string
		Port    string
		Version string
	}

	JWT struct {
		Secret string
	}
}
