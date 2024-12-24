package version

import "fmt"

type Config struct {
	Usage string `arg:"help"`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Prints Gozer's version information.`,
	}
}

func Version(cfg *Config) bool {
	// TODO: Implement
	fmt.Println(`Version is not implemented yet.`)
	fmt.Printf("\tConfig was %#v\n", cfg)
	return false
}
