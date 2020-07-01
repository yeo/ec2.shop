# Ec2.Shop

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

# Why

AWS pricing page is very slow, sometime just timing out say "Fail to
load price". I know similar service like https://ec2instances.info/ but
it's also slow and didn't have curl-able interface.

All I want is a way to compare/check price right from terminal. The URL
need to be short and easy to remember, thus `ec2.shop`.

# How accurate is the price

It's very accurate for on-demand instance, as accurate as whatever on
this page: https://aws.amazon.com/ec2/pricing/on-demand/

For spot instances, it maybe a bit outdated. The price is refresh every
5 minutes from this page: https://aws.amazon.com/ec2/spot/pricing/

# Will you maitenance this?

I need it myself and it's very cheap to keep it running given a majority
of request are cached at Cloudflare.

Otherwise, you can run it yourself. I had Dockerfile, k8s, makefile to
help you run it.
