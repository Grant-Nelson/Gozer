package driver

import (
	"context"
	"fmt"
	"io"
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
	logs     io.Writer
	roots    []*packages.Package
	packages []*packages.Package
}

func New(ctx context.Context, patterns []string, request *packages.DriverRequest, logs io.Writer) Driver {
	return &driver{
		ctx:      ctx,
		patterns: patterns,
		request:  request,
		response: &packages.DriverResponse{},
		logs:     logs,
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

	// Prevent this call to packages.Load from starting up this driver otherwise
	// it will recursively start driver processes until the OS is restarted.
	env := slices.DeleteFunc(d.request.Env,
		func(e string) bool { return strings.HasPrefix(e, `GOPACKAGES`) })

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
	d.response.Roots = make([]string, len(d.roots))
	for i, root := range d.roots {
		d.response.Roots[i] = root.ID
	}

	return nil
}

func (d *driver) flattenForest() {
	d.packages = slices.Collect(packages.Postorder(d.roots))
	d.response.Packages = d.packages
}

func (d *driver) prepareResponse() {

}
