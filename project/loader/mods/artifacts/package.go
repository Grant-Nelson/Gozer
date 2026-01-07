package artifacts

import "go/token"

type Package struct {
	name        string
	path        string
	isTest      bool
	isXTest     bool
	tempFileSet *token.FileSet
}

func NewPackage(name, path string, isTest, isXTest bool, tempFileSet *token.FileSet) *Package {
	return &Package{
		name:        name,
		path:        path,
		isTest:      isTest,
		isXTest:     isXTest,
		tempFileSet: tempFileSet,
	}
}

// NewPackageForFile creates a new package for the given file.
//
// This will not change the package on the file.
func NewPackageForFile(f *File) *Package {
	return NewPackage(f.PackageName(), f.PackagePath(), f.IsTest(), f.IsXTest(), f.TempFileSet)
}

func (p *Package) Name() string { return p.name }
func (p *Package) Path() string { return p.path }

// IsTest indicates this package is a partial package containing all `IsTest` files.
func (p *Package) IsTest() bool { return p.isTest }

// IsXTest indicates this package ia a package containing all `IsXTest` files.
func (p *Package) IsXTest() bool { return p.isXTest }

func (p *Package) TempFileSet() *token.FileSet { return p.tempFileSet }

// Key gets the key for a package based on the package path and test flags.
func (p *Package) Key() string {
	switch {
	case p.IsXTest():
		return p.Path() + `#_XTest`
	case p.IsTest():
		return p.Path() + `#_Test`
	default:
		return p.Path()
	}
}
