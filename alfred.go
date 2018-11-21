package alfred

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wzulfikar/alfred/contracts"
)

type Alfred struct {
	Finders          map[string]contracts.Finder
	ResolveFindersFn func(alfred *Alfred, query *string) *map[string]contracts.Finder
}

func New(finders *[]contracts.Finder, resolveFindersFn func(alfred *Alfred, query *string) *map[string]contracts.Finder) (*Alfred, error) {
	findersMap := make(map[string]contracts.Finder, len(*finders))
	var errs []error
	for _, finder := range *finders {
		if err := finder.Init(); err != nil {
			errs = append(errs, err)
			continue
		}
		findersMap[finder.FinderName()] = finder
	}

	if len(errs) > 0 {
		err := fmt.Sprintf("found %d error(s) when initializing alfred:\n%v", len(errs), errs)
		return nil, errors.New(err)
	}

	return &Alfred{
		Finders:          findersMap,
		ResolveFindersFn: resolveFindersFn,
	}, nil
}

func (alfred *Alfred) GetFinders(finderNames ...string) *map[string]contracts.Finder {
	if len(finderNames) == 0 {
		return &alfred.Finders
	}

	finders := make(map[string]contracts.Finder)
	for _, finderName := range finderNames {
		for name, finder := range alfred.Finders {
			if strings.HasPrefix(name, finderName) {
				finders[name] = finder
			}
		}
	}
	return &finders
}

func (alfred *Alfred) GetFindersExcept(finderNames ...string) *map[string]contracts.Finder {
	if len(finderNames) == 0 {
		return &alfred.Finders
	}

	finders := make(map[string]contracts.Finder)
	for _, finderName := range finderNames {
		for name, finder := range alfred.Finders {
			if !strings.HasPrefix(name, finderName) {
				finders[name] = finder
			}
		}
	}
	return &finders
}
