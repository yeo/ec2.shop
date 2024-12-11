#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'uri'

ts=Time.now.to_i
`mkdir -p data/msk`

[{
  name: 'msk',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/msk/USD/current/msk.json?timestamp=#{ts}"
}].each do |instance_class|
  puts "#{instance_class[:name]} downloading..."

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/msk/#{instance_class[:name]}.json", region_price_data)
end
