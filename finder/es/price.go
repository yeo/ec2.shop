package es

import (
	"github.com/yeo/ec2shop/finder/common/simpleri"
)

func Discover(r string) map[string]*simpleri.Price {
	data := simpleri.Discover(&simpleri.DiscoverRequest{
		OndemandFile: "./data/es/es-ondemand.json",
		Region:       r,
		Family:       "opensearch",
		RiPrefixPath: "./data/es/es-reservedinstance-",
		NodeTypes: []string{
			"General%20purpose",
			"Compute%20optimized",
			"Memory%20optimized",
			"Storage%20optimized",
			"OR1",
		},
	})
	return data
}
