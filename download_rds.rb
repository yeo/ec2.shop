#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'cgi'

ts=Time.now.to_i

[
  { name: 'postgresql-ondemand', url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/rds/USD/current/rds-postgresql-ondemand.json" },
  { name: 'postgresql-reserved-instances-plan', url: "https://c0.b0.p.awsstatic.com/configurations/aws/rds/postgresql-reserved-instances-plan.json", }
].each do |instance_class|
  puts "fetch #{instance_class[:name]}"

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/rds/#{instance_class[:name]}.json", region_price_data)
end
