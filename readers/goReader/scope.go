package goReader

import (
	"fmt"

	"github.com/grant-nelson/Gozer/constructs/expressions"
	"github.com/grant-nelson/Gozer/constructs/types"
)

// Scope is the set of identifiers in a program block.
type Scope struct {

	// Identifiers are the identifiers connected
	Identifiers map[string]types.Type

	// Previos is the containing scope.
	Previous *Scope
}

// NewScope creates a new scope.
func NewScope(prev *Scope) *Scope {
	return &Scope{
		Identifiers: map[string]types.Type{},
		Previous:    prev,
	}
}

// Get gets the named type with the given name.
func (scope *Scope) Get(id string) *expressions.IdentifierExp {
	if t, exists := scope.Identifiers[id]; exists {
		return expressions.Identifier(id, t)
	}
	if scope.Previous == nil {
		return nil
	}
	return scope.Previous.Get(id)
}

// Add inserts a new identifier and given type into this scope.
func (scope *Scope) Add(id string, t types.Type) *expressions.IdentifierExp {
	scope.Identifiers[id] = t
	return expressions.Identifier(id, t)
}

// AddTemp inserts a new temporary variable by a unique temporary name
// with the given type into this scope. The new name is returned.
func (scope *Scope) AddTemp(t types.Type) *expressions.IdentifierExp {
	i := 0
	for {
		temp := fmt.Sprintf("gozerTemp%d", i)
		if id := scope.Get(temp); id == nil {
			return scope.Add(temp, t)
		}
		i++
	}
}
