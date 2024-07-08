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
    { name: 'previousgen-ondemand', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/previousgen-ondemand/#{CGI.escape(region['name'])}/Linux/index.json?timestamp=#{ts}" },
    { name: 'reservedinstance-3y', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ec2-reservedinstance/3%20year/No%20Upfront/#{CGI.escape(region['name'])}/Linux/Shared/index.json?timestamp=1709867848131#{ts}" },
    { name: 'reservedinstance-1y', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ec2-reservedinstance/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/Linux/Shared/index.json?timestamp=1709867848131#{ts}" },
    { name: 'reservedinstance-convertible-1y', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ec2-reservedinstance-convertible/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/Linux/Shared/index.json?timestamp=1709877151240" },
    { name: 'reservedinstance-convertible-3y', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/ec2/USD/current/ec2-reservedinstance-convertible/3%20year/No%20Upfront/#{CGI.escape(region['name'])}/Linux/Shared/index.json?timestamp=1709877151240" },
  ].each do |instance_class|
    puts "fetch #{instance_class[:name]}"

    region_price_data = HTTPX.get(instance_class[:url])
    File.write("./data/ec2/#{region['code']}-#{instance_class[:name]}.json", region_price_data)
  end
end
