#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'cgi'

ts=Time.now.to_i

[{
  name: 'rds-postgresql-ondemand',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-ondemand.json?timestamp=#{ts}"
}, {
  name: 'rds-mariadb-ondemand',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mariadb-ondemand.json?timestamp={ts}"
},{
  name: 'rds-mysql-ondemand',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mysql-ondemand.json?timestamp={ts}"
}].each do |instance_class|
  puts "#{instance_class[:name]} downloading..."

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/rds/#{instance_class[:name]}.json", region_price_data)
end


regions=JSON.parse(HTTPX.get("https://b0.p.awsstatic.com/locations/1.0/aws/current/locations.json?timestamp=#{ts}"))

puts "found #{regions.length} regions"

regions.each do |key, region|
  puts "[postgres] fetching #{region['name']}"

  [{
    name: 'rds-mysql-reservedinstance-multi-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mysql-reservedinstance/Multi-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mysql-reservedinstance-multi-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mysql-reservedinstance/Multi-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mysql-reservedinstance-single-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mysql-reservedinstance/Single-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mysql-reservedinstance-single-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mysql-reservedinstance/Single-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mariadb-reservedinstance-multi-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mariadb-reservedinstance/Multi-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mariadb-reservedinstance-multi-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mariadb-reservedinstance/Multi-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mariadb-reservedinstance-single-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mariadb-reservedinstance/Single-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-mariadb-reservedinstance-single-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-mariadb-reservedinstance/Single-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-postgresql-reservedinstance-multi-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-reservedinstance/Multi-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-postgresql-reservedinstance-multi-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-reservedinstance/Multi-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  }, {
    name: 'rds-postgresql-reservedinstance-single-az-1y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-reservedinstance/Single-AZ/1%20year/No%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  },
  {
    name: 'rds-postgresql-reservedinstance-single-az-3y',
    url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-reservedinstance/Single-AZ/3%20year/Partial%20Upfront/#{CGI.escape(region['name'])}/index.json?timestamp=#{ts}"
  }].each do |instance_class|
    puts "fetch #{instance_class[:name]} in region #{region['name']}"

    region_price_data = HTTPX.get(instance_class[:url])
    File.write("data/rds/#{region["code"]}-#{instance_class[:name]}.json", region_price_data)
  end
end
