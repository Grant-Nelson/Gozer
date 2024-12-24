package builder

import "fmt"

type Config struct {
	Usage    string   `arg:"help"`
	Lang     string   `arg:"flag, l|lang|language, The language to transpile into."`
	Verbose  bool     `arg:"flag, v|verbose, Indicates status information should be printed while building."`
	Output   string   `arg:"flag, o|out, The directory to write the resulting application out to."`
	Patterns []string `arg:"pos, patterns, One or more patterns for the root files for a project."`
	Test     bool     `arg:"skip"`
}

func DefaultConfig(defaultLang string) *Config {
	return &Config{
		Usage: `Builds a project into a specific language at a specific` +
			`output location. The outputted code may be a package, library, ` +
			`or full application that may be used in other projects.`,
		Lang:     defaultLang,
		Verbose:  false,
		Output:   `./out`,
		Patterns: []string{},
		Test:     false,
	}
}

func Build(cfg *Config) bool {
	// TODO: Implement
	fmt.Println(`Build is not implemented yet.`)
	fmt.Printf("\tConfig was %#v\n", cfg)
	return false
}
