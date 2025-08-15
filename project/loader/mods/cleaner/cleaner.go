package cleaner

type Cleaner struct{}

func New() *Cleaner {
	return &Cleaner{}
}

// TODO: Implement a modifier that cleans up the AST before it is written.
