package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestCommonCrawlDotOrg(t *testing.T) {
	domain := "bing.com"
	source := CommonCrawlDotOrg{}
	results := []interface{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	for result := range core.UniqResults(source.ProcessDomain(ctx, domain)) {
		results = append(results, result.Success)
	}

	fmt.Println(results)

	if !(len(results) >= 3) {
		t.Errorf("expected at least 3 result(s), got '%v'", len(results))
	}
}

func TestCommonCrawlDotOrgRecursive(t *testing.T) {
	domain := "bing.com"
	source := &CommonCrawlDotOrg{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(ctx, domain, options) {
		results = append(results, result)
		fmt.Println(result)

	}

	if !(len(results) >= 3) {
		t.Errorf("expected at least 3 result(s), got '%v'", len(results))
		t.Error(ctx.Err())
	}
}

// func TestCommonCrawlDotOrg_multi_threaded(t *testing.T) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CommonCrawlDotOrg{}
// 	results := []*core.Result{}
//
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for _, domain := range domains {
// 		wg.Add(1)
// 		go func(domain string) {
// 			defer wg.Done()
// 			for result := range source.ProcessDomain(domain) {
// 				t.Log(result)
// 				mx.Lock()
// 				results = append(results, result)
// 				mx.Unlock()
// 			}
// 		}(domain)
// 	}
//
// 	wg.Wait() // collect results
//
// 	if len(results) < 40 {
// 		t.Errorf("expected more than 40 results, got '%v'", len(results))
// 	}
// }
//
// func ExampleCommonCrawlDotOrg() {
// 	domain := "bing.com"
// 	source := CommonCrawlDotOrg{}
// 	results := []*core.Result{}
//
// 	for result := range source.ProcessDomain(domain) {
// 		results = append(results, result)
// 	}
//
// 	fmt.Println(len(results) >= 3)
// 	// Output: true
// }
//
// func ExampleCommonCrawlDotOrg_multi_threaded() {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CommonCrawlDotOrg{}
// 	results := []*core.Result{}
//
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for _, domain := range domains {
// 		wg.Add(1)
// 		go func(domain string) {
// 			defer wg.Done()
// 			for result := range source.ProcessDomain(domain) {
// 				mx.Lock()
// 				results = append(results, result)
// 				mx.Unlock()
// 			}
// 		}(domain)
// 	}
//
// 	wg.Wait() // collect results
//
// 	fmt.Println(len(results) >= 40)
// 	// Output: true
// }
//
// func BenchmarkCommonCrawlDotOrgSingleThreaded(b *testing.B) {
// 	domain := "google.com"
// 	source := CommonCrawlDotOrg{}
//
// 	for n := 0; n < b.N; n++ {
// 		results := []*core.Result{}
// 		for result := range source.ProcessDomain(domain) {
// 			results = append(results, result)
// 		}
// 	}
// }
//
// func BenchmarkCommonCrawlDotOrgMultiThreaded(b *testing.B) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CommonCrawlDotOrg{}
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for n := 0; n < b.N; n++ {
// 		results := []*core.Result{}
//
// 		for _, domain := range domains {
// 			wg.Add(1)
// 			go func(domain string) {
// 				defer wg.Done()
// 				for result := range source.ProcessDomain(domain) {
// 					mx.Lock()
// 					results = append(results, result)
// 					mx.Unlock()
// 				}
// 			}(domain)
// 		}
//
// 		wg.Wait() // collect results
// 	}
// }
