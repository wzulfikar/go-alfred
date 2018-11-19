package alfred

import (
	"log"
	"time"

	"github.com/wzulfikar/alfred/contracts"
)

type FinderChan struct {
	finderName  string
	itemsFound  int
	elapsedTime string
}

func Find(query string, maxQueryResult int, finders *[]contracts.Finder) (*[]contracts.Result, error) {
	now := time.Now()

	finderChan := make(chan *FinderChan, len(*finders))
	resultChan := make(chan contracts.Result)

	for i, finder := range *finders {
		go func(i int, finder contracts.Finder, query string, resultChan chan contracts.Result) {
			finderName := finder.FinderName()
			log.Printf("fetching results from finder '%s'..", finderName)

			start := time.Now()

			result, err := finder.Find(query)
			if err != nil {
				log.Printf(`error fetching data from finder "%s": %s\n`, finderName, err)
				return
			}

			for _, item := range *result {
				resultChan <- item
			}

			finderChan <- &FinderChan{finderName, len(*result), time.Since(start).String()}
		}(i, finder, query, resultChan)
	}

	countFinder := 0
	combinedResults := &[]contracts.Result{}

	// drain the resultChan
	for {
		if len(*combinedResults) >= maxQueryResult {
			log.Printf("omitting query results due to limit (capped at %d items)\n", maxQueryResult)
			break
		}
		if countFinder == len(*finders) {
			break
		}

		// synchronize results
		select {
		case result := <-resultChan:
			if (result).Description == "" {
				(result).Description = "(No description)"
			}
			*combinedResults = append(*combinedResults, result)
		case v := <-finderChan:
			countFinder++
			log.Printf(`finder "%s" found %d items in %s. remaining finders: %d`,
				v.finderName,
				v.itemsFound,
				v.elapsedTime,
				len(*finders)-countFinder)
		default:
		}
	}

	// display "Not found" indicator when there's no item found
	if len(*combinedResults) == 0 {
		noResultIdentifier := contracts.Result{
			ID:          "0",
			Title:       "Not found",
			Description: "Whoops! Your query doesn't return anything",
			Text:        query,
		}
		*combinedResults = append(*combinedResults, noResultIdentifier)
	}

	log.Printf("[DONE] queried \"%s\" across %d finders. results found: %d. time elapsed: %s",
		query,
		len(*finders),
		len(*combinedResults),
		time.Since(now))

	return combinedResults, nil
}
