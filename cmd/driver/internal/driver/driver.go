package driver

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Driver interface {
	Drive() (*packages.DriverResponse, error)
}

type driver struct {
	ctx      context.Context
	patterns []string
	request  *packages.DriverRequest
	response *packages.DriverResponse
	roots    []*packages.Package
	packs    []*packages.Package
}

func New(ctx context.Context, patterns []string, request *packages.DriverRequest) Driver {
	return &driver{
		ctx:      ctx,
		patterns: patterns,
		request:  request,
		response: &packages.DriverResponse{},
	}
}

func (d *driver) Drive() (*packages.DriverResponse, error) {
	if err := d.loadPackages(); err != nil {
		return nil, err
	}
	d.flattenForest()

	// TODO: Finish

	d.prepareResponse()
	return d.response, nil
}

func (d *driver) loadPackages() error {
	const allNeeds = packages.NeedName |
		packages.NeedFiles |
		packages.NeedImports |
		packages.NeedDeps

	env := slices.DeleteFunc(d.request.Env,
		func(e string) bool { return strings.HasPrefix(e, `GOPACKAGES`) })
	fmt.Fprintf(os.Stderr, "Env: %v\n\n", env)

	c := &packages.Config{
		Mode:       allNeeds,
		Env:        env,
		BuildFlags: d.request.BuildFlags,
		Tests:      d.request.Tests,
		Overlay:    d.request.Overlay,
	}

	roots, err := packages.Load(c, d.patterns...)
	if err != nil {
		return fmt.Errorf(`Failed to load list of file paths: %w`, err)
	}
	d.roots = roots
	return nil
}

func (d *driver) flattenForest() {
	d.packs = slices.Collect(packages.Postorder(d.roots))
	for i, pkg := range d.packs {
		fmt.Fprintf(os.Stderr, "(%d) ID = %q\n", i, pkg.ID)
	}
}

func (d *driver) prepareResponse() {

	//d.response.Roots

}
