package common

// There are service that the cpu and memory data isn't available in the JSON
// feed but as a part of the price template itself,
// Example: view-source:https://aws.amazon.com/amazon-mq/pricing/
//
// Those rarely change so instead we hardcode a common mapping
var InstanceToAttb = map[string]*PriceAttribute{
	"t3.micro": &PriceAttribute{
		VCPU:      2,
		MemoryGib: 1,

		Memory: "2 GiB",
	},
	"m5.large": &PriceAttribute{
		VCPU:      2,
		MemoryGib: 8,
		Memory:    "8 GiB",
	},
	"m5.xlarge": &PriceAttribute{
		VCPU:      4,
		MemoryGib: 16,
		Memory:    "16 GiB",
	},
	"m5.2xlarge": &PriceAttribute{
		VCPU:      8,
		MemoryGib: 32,
		Memory:    "32 GiB",
	},
	"m5.4xlarge": &PriceAttribute{
		VCPU:      16,
		MemoryGib: 64,
		Memory:    "64 GiB",
	},

	"t2.micro": &PriceAttribute{
		VCPU:      1,
		MemoryGib: 6,
		Memory:    "1 GiB",
	},

	"m4.large": &PriceAttribute{
		VCPU:      2,
		MemoryGib: 8,
		Memory:    "8 GiB",
	},
}
