package alfred

import (
	"log"
	"sync"

	"github.com/wzulfikar/alfred/contracts"
)

func Find(query string, maxQueryResult int, finders *[]contracts.Finder) (*[]contracts.Result, error) {
	var wg sync.WaitGroup
	wg.Add(len(*finders))

	results := make([][]contracts.Result, len(*finders))
	for i, finder := range *finders {
		go func(i int, finder contracts.Finder, query string, results *[][]contracts.Result, wg *sync.WaitGroup) {
			defer wg.Done()
			finderName := finder.FinderName()
			result, err := finder.Find(query)
			if err != nil {
				log.Printf("error fetching data from finder '%s': %s\n", finderName, err)
				return
			}
			log.Printf("%d data fetched from finder '%s'\n", len(*result), finderName)
			(*results)[i] = *result
		}(i, finder, query, &results, &wg)
	}
	wg.Wait()

	log.Printf("finished fetching results from %d finders\n", len(*finders))

	counter := 0
	combinedResults := &[]contracts.Result{}
	for _, driverResult := range results {
		for _, item := range driverResult {
			if item.Description == "" {
				item.Description = "(No description)"
			}
			*combinedResults = append(*combinedResults, item)
			counter++

			if counter >= maxQueryResult {
				log.Printf("omitting query results due to limit (capped at %d items)\n", maxQueryResult)
				break
			}
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

	return combinedResults, nil
}
