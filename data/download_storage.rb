#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'cgi'

ts=Time.now.to_i
`mkdir -p data/ebs`
`mkdir -p data/efs`

[{
  name: 'ebs',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ebs.json?timestamp=#{ts}",

  name: 'efs',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/efs/USD/current/efs.json?timestamp=#{ts}"
},].each do |instance_class|
  puts "#{instance_class[:name]} downloading..."

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/#{name}/#{name}.json", region_price_data)
end
