package config

type Main struct {
	Usage   string   `arg:"help"`
	Lang    string   `arg:"flag,l|lang|language,The language to output as."`
	Build   *Build   `arg:"tool,b|build,Builds a project."`
	Run     *Run     `arg:"tool,r|run,Builds a project then runs it."`
	Test    *Test    `arg:"tool,t|test,Builds and tests all or part of a project"`
	Version *Version `arg:"tool,v|version,Shows the version."`
}

type Build struct {
	Usage    string   `arg:"help"`
	Verbose  bool     `arg:"flag,v|verbose,Indicates status information should be printed while building."`
	Patterns []string `arg:"pos,patterns,One or more patterns for the root files for a project."`
}

type Run struct {
	Usage    string   `arg:"help"`
	Verbose  bool     `arg:"flag,v|verbose,Indicates status information should be printed while building and running."`
	Patterns []string `arg:"pos,patterns,One or more patterns for the root files for a project."`
}

type Test struct {
	Usage    string   `arg:"help"`
	Verbose  bool     `arg:"flag,v|verbose,Indicates status information should be printed while building and testing."`
	Patterns []string `arg:"pos,patterns,One or more patterns for the root files for a project."`
}

type Version struct {
	Usage string `arg:"help"`
}
