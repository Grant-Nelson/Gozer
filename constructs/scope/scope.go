package scope

type Scope int

const (
	Private Scope = iota
	Public
	Internal
	Protected
)

func (s Scope) String() string {
	switch s {
	case Private:
		return `private`
	case Public:
		return `public`
	case Internal:
		return `internal`
	case Protected:
		return `protected`
	default:
		return `unknown scope`
	}
}
