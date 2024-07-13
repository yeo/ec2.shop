package rds

import (
	"testing"

	"github.com/gookit/goutil/dump"

	"github.com/yeo/ec2shop/finder/common"
)

func TestLoadPriceFromJSON(t *testing.T) {
	common.LoadRegions()
	prices := Discover("us-east-1")
	dump.P(prices)
}
