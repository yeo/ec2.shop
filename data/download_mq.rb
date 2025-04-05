#!/usr/bin/env ruby3

require 'httpx'
require 'json'

ts=Time.now.to_i
`mkdir -p data/mq`

[{
  name: 'mq',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/mq/USD/current/mq.json?timestamp=#{ts}"
}].each do |instance_class|
  puts "#{instance_class[:name]} downloading..."

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/mq/#{instance_class[:name]}.json", region_price_data)
end
