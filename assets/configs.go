package assets

import "embed"

var (
	//go:embed _configs/*.toml
	configDir embed.FS
)

func LoadConfig(filename string) ([]byte, error) {
	return configDir.ReadFile("_configs/" + filename)
}
