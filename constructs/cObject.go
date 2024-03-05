package constructs

type CObject interface {
	Named
	Methods() CMethodSet
}
