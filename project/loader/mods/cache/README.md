# Cache Modifier

## Design Decisions

### About Imports

Since packages should be able to be build separately and linked when first
starting up, it may be possible that when an import has updated (even have
different exports) and the importing package won't also need updating.
This assumes that if the import changed enough that, the importing package
would also have to update. If the importing package didn't update, then the
import change was none disruptive.

To that end, if a cache is out-of-date, assume it can be rebuilt without
invalidating the cache of any package that imports it. If this is not true
for a target language, add a setting to allow invalidation of caches based
on dependency tree.

## TODO

- Need to add a manifest file to manage several different builds for the
  same package built for different build flags and versions. The build flags
  should be filtered to flags that actually affect the package.

- Need to check that the cache is newer than all the source files and
  embedded files.
  Determine if there is a way to determine the augmenter files and other
  things that may invalidate the cache, or add to the developer notes
  that when modifying the augmenter files, the cache has to be disabled.

- Check that the transpiled code exists, otherwise we should skip the
  cache so that the files are loaded and the code can be transpiled.

- Also note that the stored cache is put into a temporary location
  until the resulting transpiled code has been written, then the stored cache
  and transpiled code can be moved to the location caches are read from.
