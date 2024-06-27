# https://ec2.shop

Get ec2 price right from your terminal

```
curl 'https://ec2.shop'
```

If you prefer json, use:

```
curl -H 'accept: json' 'https://ec2.shop'
```

If you want to search for a certain instance:

```
curl 'https://ec2.shop?filter=i3'
curl 'https://ec2.shop?filter=ssd'
```

# Advanced Filter

We also support expression in filter so you can do comparison like this

```
curl 'https://ec2.shop?filter=ssd,mem>=32,mem<=64,cpu>=2,cpu<=4'
```

The pharse `ssd,mem>=32,mem<=64,cpu>=2,cpu<=4` can also be entered into our
search box to filter the desire instance.

We support below field:

- **mem**: filter based on memory in Gib
- **cpu**: filter based on cpu
- **price**: hourly price
- **spotprice**: hourly spot price

# Why

AWS pricing page is very slow, sometime just timing out say "Fail to
load price". I know similar service like https://ec2instances.info/ but
it's also slow and didn't have curl-able interface.

All I want is a way to compare/check price right from terminal. The URL
need to be short and easy to remember, thus `ec2.shop`.

# How accurate is the price

It's very accurate for on-demand instance, as accurate as whatever on
this page: https://aws.amazon.com/ec2/pricing/on-demand/

For spot instances, The price is refresh every 2.5 minutes from this page: https://aws.amazon.com/ec2/spot/pricing/ 

The spot instances price may change in 5 minutes, so we migh have a slightly outdate but given our fetch schedule(twice per 5 minutes) I think we're pretty good there.


# Will you maitenance this?

I need it myself and it's very cheap to keep it running mngiven a majority
of request are cached at Cloudflare.

Otherwise, you can run it yourself. I had Dockerfile, k8s, makefile to
help you run it.

## API Document

We support either text base or json base response. text base is useful
in text processing with `awk`. text base is default mode. To use JSON,
simply pass a `accept: json` header.

```
curl -H 'accept: json' 'https://ec2.shop'
```

To filter out response result:

Example, to find all `*.large` instance type:

```
curl -H 'accept: json' 'https://ec2.shop?filter=.large'
```

To find instance support ssd:

```
curl -H 'accept: json' 'https://ec2.shop?filter=ssd'
```

The filter parameter is an `or` query type, so you can do this:

```
curl 'https://ec2.shop?filter=t2.medium,t3.medium'
```

The text response looks like this:

```
Instance Type    Memory             vCPUs  Storage               Network             Price       Monthly     Spot Price
c5d.9xlarge      72 GiB          36 vCPUs  1 x 900 NVMe SSD      10 Gigabit          1.7280      1261.440    0.7175
m5dn.24xlarge    384 GiB         96 vCPUs  4 x 900 NVMe SSD      100 Gigabit         6.5280      4765.440    1.6323
m6g.large        8 GiB            2 vCPUs  EBS only              Up to 10 Gigabit    0.0770      56.210      0.0357
m5.xlarge        16 GiB           4 vCPUs  EBS only              Up to 10 Gigabit    0.1920      140.160     0.0806
a1.metal         32 GiB          16 vCPUs  EBS only              Up to 10 Gigabit    0.4080      297.840     0.1343
```

All price are for Linux instance. For JSON, the response contains these:

```
    {
      "InstanceType": "r3.xlarge",
      "Memory": "30.5 GiB",
      "VCPUS": 4,
      "Storage": "1 x 80 SSD",
      "Network": "Moderate",
      "Cost": 0.333,
      "MonthlyPrice": 243.09,
      "SpotPrice": "0.0650"
    }
```

Unfortunately the `SpotPrice` is a string :-( because sometime it
contains this text: `"SpotPrice": "NA"` when that instance type isn't
available for purchase on Spot Instance(as in, they are only available
for on-demand).

# Icon

Use price icon by https://www.iconfinder.com/WTicon

