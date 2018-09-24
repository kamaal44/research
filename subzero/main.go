package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/subfinder/research/core"
	"github.com/subfinder/research/core/sources"
)

var sourcesList = []core.Source{
	&sources.ArchiveIs{},
	&sources.CertSpotter{},
	&sources.CommonCrawlDotOrg{},
	&sources.CrtSh{},
	&sources.FindSubdomainsDotCom{},
	&sources.HackerTarget{},
	&sources.Riddler{},
	&sources.Threatminer{},
	&sources.WaybackArchive{},
	&sources.DNSDbDotCom{},
	&sources.Bing{},
	&sources.Yahoo{},
	&sources.Baidu{},
	&sources.Entrust{},
	&sources.ThreatCrowd{},
}

func main() {
	results := make(chan *core.Result)
	jobs := sync.WaitGroup{}
	var cmdEnumerateVerboseOpt bool
	var cmdEnumerateInsecureOpt bool
	var cmdEnumerateLimitOpt int

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := &core.EnumerationOptions{
		Sources: sourcesList,
		Context: ctx,
		Cancel:  cancel,
	}

	var cmdEnumerate = &cobra.Command{
		Use:   "enumerate [domains to enumerate]",
		Short: "Enumerate subdomains for the given domains.",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if cmdEnumerateInsecureOpt {
				sourcesList = append(sourcesList, &sources.PTRArchiveDotCom{})
				sourcesList = append(sourcesList, &sources.DogPile{})
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, domain := range args {
				jobs.Add(1)
				go func(domain string) {
					defer jobs.Done()
					for result := range core.EnumerateSubdomains(domain, opts) {
						results <- result
					}
				}(domain)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			var count = 0
			for result := range results {
				count++
				if result.IsSuccess() {
					fmt.Println(result.Type, result.Success)
				} else if cmdEnumerateVerboseOpt {
					fmt.Println(result.Type, result.Failure)
				}
				if cmdEnumerateLimitOpt != 0 && cmdEnumerateLimitOpt == count {
					cancel()
				}
			}
		},
	}
	cmdEnumerate.Flags().IntVar(&cmdEnumerateLimitOpt, "limit", 0, "Limit the reported results to the given number.")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateVerboseOpt, "verbose", false, "Show errors and other available diagnostic information.")
	cmdEnumerate.Flags().BoolVar(&cmdEnumerateInsecureOpt, "insecure", false, "Use potentially insecure sources using http.")

	var rootCmd = &cobra.Command{Use: "subzero"}
	rootCmd.AddCommand(cmdEnumerate)
	rootCmd.Execute()

	jobs.Wait()
}
