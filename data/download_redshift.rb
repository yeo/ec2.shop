#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'uri'

ts=Time.now.to_i
`mkdir -p data/redshift`

[{
  name: 'redshift',
  url: "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/redshift/USD/current/redshift.json?timestamp=1720953273260?timestamp=#{ts}"
}].each do |instance_class|
  puts "#{instance_class[:name]} downloading..."

  region_price_data = HTTPX.get(instance_class[:url])
  File.write("data/redshift/#{instance_class[:name]}.json", region_price_data)
end


regions=JSON.parse(HTTPX.get("https://b0.p.awsstatic.com/locations/1.0/aws/current/locations.json?timestamp=#{ts}"))

puts "found #{regions.length} regions"

regions.each do |key, region|
  rn = CGI.escape(region['name'])
  rc = CGI.escape(region['code'])

  puts "[redshift] fetching #{region['name']}"

  %w(1%20year 3%20year).each do |y|
    %w(No%20Upfront Partial%20Upfront All%20Upfront).each do |p|
      %w(Yes No).each do |t|
        url = "https://b0.p.awsstatic.com/pricing/2.0/meteredUnitMaps/redshift/USD/current/redshift-reservedinstance/#{rn}/#{y}/#{p}/#{t}/index.json?timestamp=#{ts}"
        puts "fetch #{url}"
        region_price_data = HTTPX.get url

        File.write("data/redshift/redshift-reservedinstance-#{y}-#{p}-#{rc}-#{t}.json", region_price_data)
      end
    end
  end
end
