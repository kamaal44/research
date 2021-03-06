package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// HackerTarget is a source to process subdomains from https://hackertarget.com
type HackerTarget struct {
	APIKey string
	lock   *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *HackerTarget) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		// get response from the API, optionally with an API key
		var resp *http.Response

		// http req
		var req *http.Request
		var err error

		// check API key
		if source.APIKey != "" {
			req, err = http.NewRequest(http.MethodGet, "https://api.hackertarget.com/hostsearch/?q="+domain+"&apikey="+source.APIKey, nil)
		} else {
			req, err = http.NewRequest(http.MethodGet, "https://api.hackertarget.com/hostsearch/?q="+domain, nil)
		}
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, nil, err))
			return
		}

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		resp, err = core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, nil, errors.New(resp.Status)))
			return
		}

		// TODO: investigate io.LimitedReader
		// read response body, extracting subdomains
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "API count exceeded - Increase Quota with Membership") {
				err := errors.New(scanner.Text())
				sendResultWithContext(ctx, results, core.NewResult("hackertarget", nil, err))
				return
			}
			str := domainExtractor([]byte(strings.Split(scanner.Text(), ",")[0]))
			if str != "" {
				if !sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(hackertargetLabel, nil, err))
			return
		}
	}(domain, results)
	return results
}
