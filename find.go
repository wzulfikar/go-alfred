package alfred

import (
	"log"
	"time"

	"github.com/wzulfikar/alfred/contracts"
)

var SecondsToTimeout time.Duration = 2

type FinderChan struct {
	finderName  string
	itemsFound  int
	elapsedTime string
}

func (alfred *Alfred) Find(query string, maxQueryResult int) (*[]contracts.Result, error) {
	resultIds := make(map[string]bool)

	var finders map[string]*contracts.Finder
	if alfred.ResolveFindersFn == nil {
		finders = alfred.Finders
	} else {
		finders = alfred.ResolveFindersFn(alfred, &query)
	}

	log.Printf("executing `alfred.Find()` with %ds timeout. using %d finders: %v",
		SecondsToTimeout,
		len(finders),
		finders)

	now := time.Now()

	timeoutChan := make(chan bool, 1)
	finderChan := make(chan *FinderChan, len(finders))
	resultChan := make(chan contracts.Result)

	go func() {
		time.Sleep(SecondsToTimeout * time.Second)
		timeoutChan <- true
	}()

	for finderName, finder := range finders {
		go func(finderName string, finder contracts.Finder, query string, resultChan chan contracts.Result) {
			start := time.Now()

			result, err := finder.Find(query)
			if err != nil {
				log.Printf("[ERROR] fetching data from finder \"%s\": %s\n", finderName, err)
				finderChan <- &FinderChan{finderName, 0, time.Since(start).String()}
				return
			}

			for _, item := range *result {
				resultChan <- item
			}

			finderChan <- &FinderChan{finderName, len(*result), time.Since(start).String()}
		}(finderName, *finder, query, resultChan)
	}

	countFinder := 0
	combinedResults := &[]contracts.Result{}
	timeoutReached := false

	// drain the resultChan
	for {
		if timeoutReached || countFinder == len(finders) {
			break
		}
		if len(*combinedResults) >= maxQueryResult {
			log.Printf("omitting query results due to limit (capped at %d items)\n", maxQueryResult)
			break
		}

		// synchronize results
		select {
		case result := <-resultChan:
			if resultIds[result.ID] {
				// log.Println("skipped duplicate result id:", result.ID)
				continue
			}

			if (result).Description == "" {
				(result).Description = "(No description)"
			}
			*combinedResults = append(*combinedResults, result)
			resultIds[result.ID] = true
		case v := <-finderChan:
			countFinder++
			log.Printf("- finder \"%s\" found %d items in %s. remaining finders: %d",
				v.finderName,
				v.itemsFound,
				v.elapsedTime,
				len(finders)-countFinder)
		case <-timeoutChan:
			timeoutReached = true
			log.Printf("! terminating finders due to %ds timeout. finders skipped: %d",
				SecondsToTimeout,
				len(finders)-countFinder)
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
		len(finders),
		len(*combinedResults),
		time.Since(now))

	return combinedResults, nil
}
