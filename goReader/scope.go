package transpiler

import (
	"fmt"

	"github.com/grant-nelson/Gozer/constructs"
)

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
func (scope *Scope) Get(id string) *constructs.IdentifierExp {
	if t, exists := scope.Identifiers[id]; exists {
		return constructs.Identifier(id, t)
	}
	if scope.Previous == nil {
		return nil
	}
	return scope.Previous.Get(id)
}

// Add inserts a new identifier and given type into this scope.
func (scope *Scope) Add(id string, t constructs.Type) *constructs.IdentifierExp {
	scope.Identifiers[id] = t
	return constructs.Identifier(id, t)
}

// AddTemp inserts a new temporary variable by a unique temporary name
// with the given type into this scope. The new name is returned.
func (scope *Scope) AddTemp(t constructs.Type) *constructs.IdentifierExp {
	i := 0
	for {
		temp := fmt.Sprintf("gozerTemp%d", i)
		if id := scope.Get(temp); id == nil {
			return scope.Add(temp, t)
		}
		i++
	}
}
