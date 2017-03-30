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
func (scope *Scope) Get(id string) (constructs.Type, bool) {
	if t, exists := scope.Identifiers[id]; exists {
		return t, true
	}
	if scope.Previous == nil {
		return nil, false
	}
	return scope.Previous.Get(id)
}

// Add inserts a new identifier and given type into this scope.
func (scope *Scope) Add(id string, t constructs.Type) {
	scope.Identifiers[id] = t
}

// AddTemp inserts a new temporary variable by a unique temporary name
// with the given type into this scope. The new name is returned.
func (scope *Scope) AddTemp(t constructs.Type) string {
	i := 0
	for {
		temp := fmt.Sprintf("temp%d", i)
		if _, exists := scope.Get(temp); !exists {
			scope.Add(temp, t)
			return temp
		}
		i++
	}
}
