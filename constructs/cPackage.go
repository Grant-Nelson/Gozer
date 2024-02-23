package constructs

type CPackage struct {
	Name    string
	Path    string
	Imports []*CImport
}

func (p *CPackage) ImportForPath(path string) *CImport {
	for _, i := range p.Imports {
		if i.Path == path {
			return i
		}
	}
	return nil
}
