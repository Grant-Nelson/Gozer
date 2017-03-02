package transpiler

import "github.com/grant-nelson/Gozer/constructs"

// Scope is the set of identifiers in a program block.
type Scope struct {

	// Identifiers are the identifiers connected
	Identifiers map[string]constructs.Type

	// Previos is the containing scope.
	Previous *Scope
}

// NewScope creates a new scope.
func NewScope(prev *Scope) *Scope {
	return &Scope{
		Identifiers: map[string]constructs.Type{},
		Previous:    prev,
	}
}

// Get gets the named type with the given name.
func (scope *Scope) Get(id string) (constructs.Type, bool) {
	if t, exists := scope.Identifiers[id]; exists {
		return t, true
	}
	if scope.Previous == nil {
		return nil, false
	}
	return scope.Previous.Get(id)
}

// Add inserts a new identifier and named type into this scope.
func (scope *Scope) Add(id string, t constructs.Type) {
	scope.Identifiers[id] = t
}
