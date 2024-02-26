package constructs

type CImport struct {
	Path    string
	Package *CPackage
}

func NewImport(path string) *CImport {
	return &CImport{
		Path:    path,
		Package: nil,
	}
}
