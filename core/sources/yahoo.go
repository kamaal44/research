package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// Yahoo is a source to process subdomains from https://yahoo.com
type Yahoo struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Yahoo) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	var resultLabel = "yahoo"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			url := "https://search.yahoo.com/search?p=site:" + domain + "&b=" + strconv.Itoa(currentPage*10) + "&pz=10&bct=0&xargs=0"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}

			req.Cancel = ctx.Done()
			req.WithContext(ctx)

			resp, err := core.HTTPClient.Do(req)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				if ctx.Err() != nil {
					resp.Body.Close()
					return
				}
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
							resp.Body.Close()
							return
						}
					}
				}
			}

			err = scanner.Err()

			if err != nil {
				resp.Body.Close()
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
