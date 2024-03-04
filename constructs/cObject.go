package constructs

type CObject interface {
	Named
	Directives() CDirectives
	Methods() CMethodSet
}
