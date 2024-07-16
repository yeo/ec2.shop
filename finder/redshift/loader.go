package redshift

import (
	"github.com/yeo/ec2shop/finder/common"
	"github.com/yeo/ec2shop/finder/common/simpleri"
)

func Discover(r string) map[string]*simpleri.Price {
	data := simpleri.Discover(&simpleri.DiscoverRequest{
		OndemandFile: "./data/redshift/redshift.json",
		Region:       r,
		Family:       "redshift",
		RiPrefixPath: "./data/redshift/redshift-reservedinstance-",
		NodeTypes: []string{
			"Yes",
			"No",
		},
		FilterFunc: func(name string, a *common.PriceAttribute) bool {
			// Redshift for some reason has the price item line where mem/cpu is
			// all blank, probably some serverless so we will exclude it
			if a.Memory == "" || a.RawVCPU == "" {
				return false
			}

			return true
		},
	})
	return data
}
