package fetcher

type Fetcher interface {
	Fetch(path string) ([]byte, error)
}
