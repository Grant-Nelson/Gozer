package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/Grant-Nelson/Gozer/avail/args"
)

var Config = &struct {
	Help     string `arg:"help"`
	BasePath string `arg:"optional,pos,BasePath,Base path to serve files from"`
	Verbose  bool   `arg:"optional,flag,v|verbose,Print verbose output"`
	Port     string `arg:"optional,flag,p|port,Address port to serve to"`
	Minify   bool   `arg:"optional,flag,m|mini,Minifies any *.js being built"`
}{
	Help: `Serve will serve file(s) at the base path. ` +
		`If a requested file is *.js that doesn't exist on disk, ` +
		`then the server will look for a *.ts with the same name. ` +
		`If a *.ts exists, then it will be compiled for the missing *.js ` +
		`and source mapping will also be generated. If index.html doesn't ` +
		`exist on disk, a very simple one will be generated to load all of ` +
		`the *.js and *.ts scripts.`,
	BasePath: `.`,
	Verbose:  false,
	Port:     `8080`,
	Minify:   false,
}

func main() {
	if !args.Parse(Config) {
		return
	}

	handle := http.FileServerFS(builderFileSystem{})
	fmt.Printf("listening on http://localhost:%s\n", Config.Port)
	log.Fatal(http.ListenAndServe(`:`+Config.Port, handle))
}

func logf(format string, v ...any) {
	if Config.Verbose {
		log.Printf(format, v...)
	}
}

type builderFileSystem struct{}

func (builderFileSystem) Open(name string) (fs.File, error) {
	logf("Request: %q\n", name)
	f, err := http.Dir(Config.BasePath).Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) && strings.HasSuffix(name, `index.html`) {
			f, err := generateFile(name)
			if err == nil {
				return f, nil
			}
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
		}
		logf("Error: %v\n", err)
		return nil, err
	}
	return f, nil
}

func generateFile(name string) (http.File, error) {
	if strings.HasSuffix(name, `index.html`) {
		return generateIndexPage()
	}

	if strings.HasSuffix(name, `.js`) {
		f, err := readAndBuild(name, `.ts`, api.LoaderTS)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		f, err = readAndBuild(name, `.tsx`, api.LoaderTSX)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		return readAndBuild(name, `.jsx`, api.LoaderJSX)
	}

	return nil, fs.ErrNotExist
}

func readAndBuild(jsName, ext string, loader api.Loader) (http.File, error) {
	name := strings.TrimSuffix(jsName, filepath.Ext(jsName)) + ext
	f, err := http.Dir(Config.BasePath).Open(name)
	if err != nil {
		return nil, err
	}

	buf := &strings.Builder{}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	content, err := buildJS(buf.String(), loader)
	if err != nil {
		return nil, err
	}

	logf("Built %q from %q\n", jsName, name)

	return &pseudoFile{
		name:    jsName,
		source:  f,
		content: strings.NewReader(content),
	}, nil
}

func buildJS(content string, loader api.Loader) (string, error) {
	options := api.TransformOptions{
		Target:      api.ES2015,
		TreeShaking: api.TreeShakingFalse,
		Sourcemap:   api.SourceMapInline,
		Loader:      loader,
	}

	if Config.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.KeepNames = true
	}

	result := api.Transform(string(content), options)
	for _, w := range result.Warnings {
		logf("Warning: %d:%d: %s\n%s\n", w.Location.Line, w.Location.Column, w.Text, w.Location.LineText)
	}
	if errCount := len(result.Errors); errCount > 0 {
		for _, e := range result.Errors {
			logf("Error: %d:%d: %s\n%s\n", e.Location.Line, e.Location.Column, e.Text, e.Location.LineText)
		}
		return ``, fmt.Errorf(`JS transform failed with %d errors.`, errCount)
	}

	return string(result.Code), nil
}

func generateIndexPage() (http.File, error) {

	// TODO: Implement

	return nil, fmt.Errorf(`not implemented`)
}

type pseudoFile struct {
	name    string
	source  http.File
	content *strings.Reader
}

func (pf *pseudoFile) Close() error {
	return pf.source.Close()
}

func (pf *pseudoFile) Read(p []byte) (n int, err error) {
	return pf.source.Read(p)
}

func (pf *pseudoFile) Seek(offset int64, whence int) (int64, error) {
	return pf.source.Seek(offset, whence)
}

func (pf *pseudoFile) Readdir(count int) ([]fs.FileInfo, error) {
	return pf.source.Readdir(count)
}

func (pf *pseudoFile) Stat() (fs.FileInfo, error) {
	stat, err := pf.source.Stat()
	if err != nil {
		return nil, err
	}
	return &pseudoStat{
		name: pf.name,
		size: int64(pf.content.Len()),
		stat: stat,
	}, nil
}

type pseudoStat struct {
	name string
	size int64
	stat fs.FileInfo
}

func (ps *pseudoStat) Name() string       { return ps.name }
func (ps *pseudoStat) Size() int64        { return ps.size }
func (ps *pseudoStat) Mode() fs.FileMode  { return ps.stat.Mode() }
func (ps *pseudoStat) ModTime() time.Time { return ps.stat.ModTime() }
func (ps *pseudoStat) IsDir() bool        { return ps.stat.IsDir() }
func (ps *pseudoStat) Sys() any           { return ps.stat.Sys() }
