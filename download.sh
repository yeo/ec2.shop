#!/usr/bin/env bash

set -xeuo pipefail

ts=`date +%s`

for region in `curl 'https://a0.p.awsstatic.com/pricing/1.0/ec2/manifest.json?timestamp=1593164829522'  | jq -r '.ec2[]'`; do

  curl "https://a0.p.awsstatic.com/pricing/1.0/ec2/region/$region/ondemand/linux/index.json?timestamp=$ts" -o data/$region-ondemand.json
done
