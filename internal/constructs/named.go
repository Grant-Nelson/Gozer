package constructs

type INamed interface {
	String() string
	Rename(name string)
	Exported() bool
	Anomalous() bool
}

type namedImp struct {
	name string
}

func (imp *namedImp) String() string {
	return imp.name
}

func (imp *namedImp) Rename(name string) {
	imp.name = name
}

func (imp *namedImp) Exported() bool {
	if len(imp.name) > 0 {
		firstLetter := imp.name[0]
		return firstLetter >= 'A' && firstLetter <= 'Z'
	}
	return false
}

func (imp *namedImp) Anomalous() bool {
	return len(imp.name) <= 0
}

func newName(name string) namedImp {
	// TODO: Check that names are valid
	return namedImp{name: name}
}
