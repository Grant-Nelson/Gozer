package transpiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/grant-nelson/Gozer/constructs/types"
	"github.com/grant-nelson/Gozer/framework"
	"github.com/grant-nelson/Gozer/msg"
)

// GoReader is a GO parser and Dart writer for transpileing.
type GoReader struct {

	// Program is the program being loaded.
	Program *types.ProgramType

	fileSet *token.FileSet
	sources map[string]*Source
	log     *msg.Logger
	testing bool
	results map[string]string
	goPath  string
	goRoot  string
}

// NewGoReader creates a new Go reader.
func NewGoReader() *GoReader {
	gr := &GoReader{
		Program: types.Program(),
		fileSet: token.NewFileSet(),
		sources: map[string]*Source{},
		log:     msg.NewLogger(),
		testing: false,
		results: map[string]string{},
		goPath:  os.Getenv("GOPATH"),
		goRoot:  os.Getenv("GOROOT"),
	}

	// Add all prebuilt packages.
	framework.BuiltinPrebuild(gr.Program)
	framework.FmtPrebuild(gr.Program)
	framework.IOPrebuild(gr.Program)
	return gr
}

// Logger gets the logger that the parser writes to.
func (gr *GoReader) Logger() *msg.Logger {
	return gr.log
}

// TestingEnabled indicates if testing is enabled.
func (gr *GoReader) TestingEnabled() bool {
	return gr.testing
}

// EnableTesting enables or disables testing for outputting
// and adding unit-testing while reading packages.
func (gr *GoReader) EnableTesting(test bool) {
	gr.testing = test
}

// AddFolder adds all the files in the given folder.
func (gr *GoReader) AddFolder(dirPath string) *types.PackageType {
	pack := gr.getOrCreatePackage(dirPath)
	fullPath := gr.getImportPath(dirPath)
	if len(fullPath) <= 0 {
		return nil
	}
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		gr.log.Error("Failed to add folder, \"", fullPath, "\": ", err)
		return nil
	}
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".go") {
			if (!gr.testing) && strings.HasSuffix(fileName, "_test.go") {
				continue
			}
			gr.addSource(pack, path.Join(fullPath, fileName), path.Join(dirPath, fileName), nil)
		}
	}
	return pack
}

// AddFile add the file at the given path.
// Returns true is added, false if path already exists.
func (gr *GoReader) AddFile(filePath string) (bool, *Source) {
	fullPath := gr.getImportPath(filePath)
	return gr.addSource(nil, fullPath, filePath, nil)
}

// AddCode adds the given code lines.
// The given path is the name to store this code under.
// Returns true is added, false if path already exists.
func (gr *GoReader) AddCode(filePath string, code ...string) (bool, *Source) {
	return gr.addSource(nil, filePath, filePath, strings.Join(code, "\n"))
}

// Transpile converts the code and writes the resulting files.
func (gr *GoReader) Transpile() {
	for _, src := range gr.sources {
		src.ProcessTypes()
	}
	for _, src := range gr.sources {
		src.ProcessBodies()
	}
	// TODO: Need to connect all interfaces across all the classes.
	// for i := gr.packages.Length() - 1; i >= 0; i-- {
	// 	pack := gr.packages.At(i).Data.(*types.PackageType)
	// 	pack.ProcessBodies()
	// }
}

// Results gets the transpiled results.
func (gr *GoReader) Results() map[string]string {
	return gr.results
}

//============================================================================

// fileExists checks if the file exists.
func (gr *GoReader) fileExists(elems ...string) (string, bool) {
	file := path.Join(elems...)
	_, err := os.Stat(file)
	return file, err == nil
}

// srcPath gets the path to the GO source directory.
func (gr *GoReader) getImportPath(path string) string {
	if file, exists := gr.fileExists(path); exists {
		return file
	}
	if len(gr.goPath) > 0 {
		if file, exists := gr.fileExists(gr.goPath, "src", path); exists {
			return file
		}
	}
	if len(gr.goRoot) > 0 {
		if file, exists := gr.fileExists(gr.goRoot, "src", path); exists {
			return file
		}
	}

	gr.log.Error("Failed to find import path: ", path)
	return ""
}

// addSource adds GO source code to the transpiler.
func (gr *GoReader) addSource(pack *types.PackageType, fullPath string, importPath string, code interface{}) (bool, *Source) {
	dirPath, fileName := path.Split(importPath)
	if strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath[:len(dirPath)-1]
	}
	if pack == nil {
		pack = gr.getOrCreatePackage(dirPath)
	}
	if source, exists := gr.sources[fileName]; exists {
		return false, source
	}

	// Create new source.
	data, err := parser.ParseFile(gr.fileSet, fullPath, code, parser.ParseComments)
	if err != nil {
		gr.log.Error("Failed to add source, ", fullPath, ": ", err)
		return false, nil
	}

	source := NewSource(gr.log, gr.fileSet)
	source.Path = importPath
	source.Package = pack
	gr.sources[fileName] = source
	source.Data = data

	// Resolve imports for source.
	gr.addImport(pack, source, "builtin", "builtin")
	for _, spec := range source.Data.Imports {
		importPath, importShort := gr.readImport(spec)
		gr.addImport(pack, source, importPath, importShort)
	}
	return true, source
}

// getOrCreatePackage gets or creates the package for the given path.
func (gr *GoReader) getOrCreatePackage(dirPath string) *types.PackageType {
	pack, exists := gr.Program.Packages[dirPath]
	if !exists {
		pack = types.Package()
		gr.Program.Packages[dirPath] = pack
	}
	return pack
}

// addImport adds an import to the given source and package.
func (gr *GoReader) addImport(pack *types.PackageType, source *Source, path string, short string) {
	importPack := gr.resolveImport(path)

	// Add import to source with short name.
	source.Imports[short] = importPack

	// Add import to package with generalized name.
	pack.Imports[path] = importPack
}

// readImport creates a new import from the given spec.
func (gr *GoReader) readImport(spec *ast.ImportSpec) (string, string) {
	path := spec.Path.Value
	if strings.HasPrefix(path, "\"") {
		path = path[1:]
	}
	if strings.HasSuffix(path, "\"") {
		path = path[:len(path)-1]
	}

	short := path
	if spec.Name != nil {
		short = spec.Name.Name
	} else if last := strings.LastIndex(path, "/"); last > 0 {
		short = path[last+1:]
	}
	return path, short
}

// resolveImport determines if the import exists, if it is part
// of the framework, or the import is added into sources.
func (gr *GoReader) resolveImport(importPath string) *types.PackageType {
	if pack, exists := gr.Program.Packages[importPath]; exists {
		return pack
	}
	return gr.AddFolder(importPath)
}
