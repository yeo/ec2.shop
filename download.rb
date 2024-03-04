#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'cgi'

ts=Time.now.to_i

regions=JSON.parse(HTTPX.get("https://b0.p.awsstatic.com/locations/1.0/aws/current/locations.json?timestamp=#{ts}"))

puts "found #{regions.length} regions"
regions.each do |key, region|
  puts "fetching #{region['name']}"

  [
    { name: 'ondemand', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ec2-ondemand-without-sec-sel/#{CGI.escape(region['name'])}/Linux/index.json?timestamp=#{ts}"},
    { name: 'previousgen-ondemand', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/previousgen-ondemand/#{CGI.escape(region['name'])}/Linux/index.json?timestamp=#{ts}" }
  ].each do |instance_class|
    region_price_data = HTTPX.get(instance_class[:url])
    File.write("./data/#{region['code']}-#{instance_class[:name]}.json", region_price_data)
  end
end
