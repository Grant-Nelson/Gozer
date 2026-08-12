# Converter

The converter will convert from Go's AST nodes into the Gozer IR nodes.

## TODO

- Need to create a Gozer type system that allows for SSA and
  remodeling to become closer to the target language.
- Need to scrape function literals and put them into the package's functions list.
- Need to scrape locally defined types and put them into the package's type list.
- Need to read more directives such as `//go:linkname`.
- Need to keep document comments with declarations.
