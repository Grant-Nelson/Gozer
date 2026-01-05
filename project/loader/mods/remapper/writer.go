package remapper

import (
	"bytes"
	"go/printer"
	"go/token"
	"io"

	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

// Write will write the modified file to the given writer.
//
// This will not use the error group and returns any errors that occurred.
func Write(f *artifacts.File, out io.Writer) (err error) {
	// TODO: Add a defer to catch panics from the print method (it can nil pointer deref).

	cfg := &printer.Config{
		Mode:     printer.TabIndent,
		Tabwidth: 4,
	}
	return cfg.Fprint(out, f.TempFileSet(), f.File)
}

// Reload will write the file to a temporary buffer and reload it
// with the given file set to normalize the file information.
func Reload(f *artifacts.File, fileSet *token.FileSet) (*artifacts.File, error) {
	buf := &bytes.Buffer{}
	if err := Write(f, buf); err != nil {
		return nil, err
	}
	return artifacts.Load(fileSet, f.FilePath(), buf.Bytes())
}
