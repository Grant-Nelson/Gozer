package composure

type Composure interface {
	Update(other Composure) Composure
}
