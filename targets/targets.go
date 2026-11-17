package targets

import (
	"errors"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/project"
	"github.com/Grant-Nelson/Gozer/targets/target"
	"github.com/Grant-Nelson/Gozer/targets/typescript"
)

const (
	DefaultLang = `typescript`

	defaultErrorGroupLimit = 10
)

var ErrTargetUnsupported = errors.New(`unsupported target language`)

// Targets returns a list of all supported target languages.
func Targets() []target.Target {
	return []target.Target{
		typescript.Target(),
	}
}

// FindTarget finds the target with given language using the targets' alias
// names. If no target is found, nil is returned.
func FindTarget(lang string) target.Target {
	for _, t := range Targets() {
		for _, a := range t.Aliases() {
			if strings.EqualFold(a, lang) {
				return t
			}
		}
	}
	return nil
}

func findTargetOrFail(lang string) target.Target {
	if target := FindTarget(lang); target != nil {
		return target
	}
	panic(faults.From(ErrTargetUnsupported).
		With(`target`, lang))
}

// List populates the project with filenames and
// packages without parsing files or checking types.
func List(cfg *target.ListConfig) (proj *project.Project, err error) {
	if cfg.ErrGroup == nil {
		cfg.ErrGroup = faults.NewErrGroup(defaultErrorGroupLimit)
	}
	defer cfg.ErrGroup.Recover(&err)
	proj, err = findTargetOrFail(cfg.Lang).List(cfg)
	cfg.ErrGroup.Add(err)
	return proj, cfg.ErrGroup.AnyOrNil()
}

// Build performs a build to the target language.
func Build(cfg *target.BuildConfig) (err error) {
	if cfg.ErrGroup == nil {
		cfg.ErrGroup = faults.NewErrGroup(defaultErrorGroupLimit)
	}
	defer cfg.ErrGroup.Recover(&err)
	cfg.ErrGroup.Add(findTargetOrFail(cfg.Lang).Build(cfg))
	return cfg.ErrGroup.AnyOrNil()
}
