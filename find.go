package alfred

import (
	"log"
	"sync"

	"github.com/wzulfikar/alfred/contracts"
)

func Find(query string, maxQueryResult int, drivers *[]contracts.Driver) (*[]contracts.Result, error) {
	var wg sync.WaitGroup
	wg.Add(len(*drivers))

	results := make([][]contracts.Result, len(*drivers))
	for i, driver := range *drivers {
		go func(i int, driver contracts.Driver, query string, results *[][]contracts.Result, wg *sync.WaitGroup) {
			defer wg.Done()
			driverName := driver.DriverName()
			result, err := driver.Find(query)
			if err != nil {
				log.Printf("error fetching data from driver '%s': %s\n", driverName, err)
				return
			}
			log.Printf("%d data fetched from driver '%s'\n", len(*result), driverName)
			(*results)[i] = *result
		}(i, driver, query, &results, &wg)
	}
	wg.Wait()

	log.Printf("finished fetching results from %d drivers\n", len(*drivers))

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
