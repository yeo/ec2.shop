package elasticache

import (
	"github.com/yeo/ec2shop/finder/common/simpleri"
)

func Discover(elasticacheFamily, r string) map[string]*simpleri.Price {
	data := simpleri.Discover(&simpleri.DiscoverRequest{
		OndemandFile: "./data/elasticache/elasticache.json",
		Region:       r,
		Family:       elasticacheFamily,
		RiPrefixPath: "./data/elasticache/elasticache-reservedinstance-",
		NodeTypes: []string{
			"Standard",
			"Network%20optimized",
			"Memory%20optimized",
		},
	})
	return data
}
