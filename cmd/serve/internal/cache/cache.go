package cache

import (
	"fmt"
	"sync"

	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fetcher"
)

func New(next fetcher.Fetcher) fetcher.Fetcher {
	return &cacheImp{
		next:       next,
		cached:     make(map[string][]byte),
		inProgress: make(map[string]chan struct{}),
	}
}

type cacheImp struct {
	next       fetcher.Fetcher
	cached     map[string][]byte
	inProgress map[string]chan struct{}
	lock       sync.Mutex
}

func (c *cacheImp) Fetch(path string) ([]byte, error) {
	return c.syncLoad(path)()
}

type syncFunc func() ([]byte, error)

func (c *cacheImp) syncLoad(path string) syncFunc {
	c.lock.Lock()
	defer c.lock.Unlock()

	if content, found := c.cached[path]; found {
		return c.alreadyCached(content)
	}

	if ch, loading := c.inProgress[path]; loading {
		return c.alreadyInProgress(path, ch)
	}

	ch := make(chan struct{})
	c.inProgress[path] = ch
	return c.startFetching(path, ch)
}

func (c *cacheImp) alreadyCached(content []byte) syncFunc {
	return func() ([]byte, error) {
		return content, nil
	}
}

func (c *cacheImp) alreadyInProgress(path string, ch chan struct{}) syncFunc {
	return func() ([]byte, error) {
		<-ch

		c.lock.Lock()
		defer c.lock.Unlock()

		if content, found := c.cached[path]; found {
			return content, nil
		}
		return nil, fmt.Errorf(`failed to find package %q in cache after waiting load`, path)
	}
}

func (c *cacheImp) startFetching(path string, ch chan struct{}) syncFunc {
	return func() ([]byte, error) {
		content, err := c.next.Fetch(path)
		if err != nil {
			return nil, err
		}

		c.lock.Lock()
		defer c.lock.Unlock()

		c.cached[path] = content
		delete(c.inProgress, path)
		close(ch)
		return content, nil
	}
}
