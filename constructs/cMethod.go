package constructs

type CMethod interface {
	Named
	Directives() CDirectives
}
