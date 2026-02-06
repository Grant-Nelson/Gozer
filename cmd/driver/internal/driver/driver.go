package driver

import (
	"context"
	"fmt"
	"slices"

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
	lists, err := d.getFileList()
	if err != nil {
		return nil, err
	}

	return d.response, nil
}

var skipEnv = map[string]bool{
	`GOPACKAGESDRIVER`: true,
	`GOPACKAGESDEBUG`:  true,
}

const allNeeds = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedDeps

func (d *driver) getFileList() ([]*packages.Package, error) {
	env := slices.DeleteFunc(d.request.Env,
		func(e string) bool { return skipEnv[e] })

	c := &packages.Config{
		Mode:       allNeeds,
		Env:        env,
		BuildFlags: d.request.BuildFlags,
		Tests:      d.request.Tests,
	}

	pkgs, err := packages.Load(c, d.patterns...)
	if err != nil {
		return nil, fmt.Errorf(`Failed to load list of file paths: %w`, err)
	}
	return pkgs, nil
}
