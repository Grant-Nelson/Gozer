package langs

import "fmt"

type Config struct {
	Usage string `arg:"help"`
}

func DefaultConfig() *Config {
	return &Config{
		Usage: `Prints the list of languages available to transpile into ` +
			`and provides additional information about those languages, ` +
			`such as what tool the language can be run and tested with.`,
	}
}

func Langs(cfg *Config) bool {
	// TODO: Implement
	fmt.Println(`ListLang is not implemented yet.`)
	fmt.Printf("\tConfig was %#v\n", cfg)
	return false
}
