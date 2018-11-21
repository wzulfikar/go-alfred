package alfred

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/wzulfikar/go-alfred/contracts"
)

type Alfred struct {
	Finders          map[string]*contracts.Finder
	ResolveFindersFn func(alfred *Alfred, query *string) map[string]*contracts.Finder
}

func New(finders *[]contracts.Finder, resolveFindersFn func(alfred *Alfred, query *string) map[string]*contracts.Finder) (*Alfred, error) {
	findersMap := make(map[string]*contracts.Finder, len(*finders))
	var errs []error
	for i := 0; i < len(*finders); i++ {
		f := (*finders)[i]
		if err := f.Init(); err != nil {
			errs = append(errs, err)
			continue
		}
		findersMap[f.FinderName()] = &f
	}

	log.Printf("alfred initialized with %d finders: %v", len(findersMap), findersMap)

	if len(errs) > 0 {
		err := fmt.Sprintf("found %d error(s) when initializing alfred:\n%v", len(errs), errs)
		return nil, errors.New(err)
	}

	return &Alfred{
		Finders:          findersMap,
		ResolveFindersFn: resolveFindersFn,
	}, nil
}

// FindersInclude returns ONLY finders that are passed as arguments
func (alfred *Alfred) FindersInclude(finderNames ...string) map[string]*contracts.Finder {
	if len(finderNames) == 0 {
		return alfred.Finders
	}

	finders := make(map[string]*contracts.Finder)
	for _, finderName := range finderNames {
		for name := range alfred.Finders {
			if strings.HasPrefix(name, finderName) {
				finders[name] = alfred.Finders[name]
			}
		}
	}
	return finders
}

// FindersExclude returns finders while excluding those passed as arguments
func (alfred *Alfred) FindersExclude(finderNames ...string) map[string]*contracts.Finder {
	if len(finderNames) == 0 {
		return alfred.Finders
	}

	finders := make(map[string]*contracts.Finder)
	for _, finderName := range finderNames {
		for name := range alfred.Finders {
			if strings.HasPrefix(name, finderName) {
				continue
			}
			finders[name] = alfred.Finders[name]
		}
	}

	return finders
}
