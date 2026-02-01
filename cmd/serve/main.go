package main

import (
	"bytes"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
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
		`exist on disk, a very simple one will be returned that will load` +
		`a "main.js" script.`,
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
		if errors.Is(err, fs.ErrNotExist) {
			f, err := generateFile(name)
			if err == nil {
				return f, nil
			}
			if !errors.Is(err, fs.ErrNotExist) {
				logf("Error generating: %v\n", err)
				return nil, err
			}
		}
		logf("Error: %v\n", err)
		return nil, err
	}
	return f, nil
}

const (
	indexPageName = `index.html`
	faviconName   = `favicon.ico`
)

func generateFile(name string) (http.File, error) {
	switch filepath.Ext(name) {
	case `.html`:
		// TODO: Change the server embedded file to just check for any base
		// name then on "not exists" continue onto generate. Also move the
		// embedded files into a dir and just embed the whole dir.
		if filepath.Base(name) == indexPageName {
			return serveEmbeddedFile(indexPageName)
		}

	case `.ico`:
		if filepath.Base(name) == faviconName {
			return serveEmbeddedFile(faviconName)
		}

	case `.js`:
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
	start := time.Now()
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

	content, err := buildJS(name, buf.String(), loader)
	if err != nil {
		return nil, err
	}

	content = postBuildProcess(content)

	logf("Built %q from %q (%v)\n", jsName, name, time.Since(start))
	return &pseudoFile{
		name:    jsName,
		size:    int64(len(content)),
		source:  f,
		content: bytes.NewReader(content),
	}, nil
}

var importReg = regexp.MustCompile(`import.*from\s*".+(\.\w+)";`)

// postBuildProcess will search for imports that import *.ts, *.tsx, and *.jsx
// files and replace the extensions to import *.js files.
func postBuildProcess(data []byte) []byte {
	var jsExt = []byte(`.js`)
	matches := importReg.FindAllSubmatchIndex(data, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		sub := matches[i]
		if len(sub) != 4 {
			logf(`match %d found with wrong number of values: %v`, i, sub)
			continue
		}
		ext := string(data[sub[2]:sub[3]])
		switch ext {
		case `.ts`, `.tsx`, `.jsx`:
			head, tail := data[:sub[2]], data[sub[3]:]
			data = append(append(head, jsExt...), tail...)
		}
	}
	return data
}

func buildJS(sourceFile, content string, loader api.Loader) ([]byte, error) {
	options := api.TransformOptions{
		Target:      api.ES2015,
		TreeShaking: api.TreeShakingFalse,
		Sourcemap:   api.SourceMapInline,
		Sourcefile:  sourceFile,
		SourceRoot:  Config.BasePath,
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
		return nil, fmt.Errorf(`JS transform failed with %d errors.`, errCount)
	}

	return result.Code, nil
}

//go:embed favicon.ico index.html
var embeddedFiles embed.FS

func serveEmbeddedFile(name string) (http.File, error) {
	return http.FS(embeddedFiles).Open(name)
}

type pseudoFile struct {
	name    string
	size    int64
	source  http.File
	content io.ReadSeeker
}

func (pf *pseudoFile) Close() error {
	return pf.source.Close()
}

func (pf *pseudoFile) Read(p []byte) (n int, err error) {
	return pf.content.Read(p)
}

func (pf *pseudoFile) Seek(offset int64, whence int) (int64, error) {
	return pf.content.Seek(offset, whence)
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
		size: pf.size,
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
