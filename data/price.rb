#!/usr/bin/env ruby3

require 'httpx'
require 'json'
require 'uri'

# First get the index file
region_raw_data = JSON.parse(HTTPX.get("https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/index.json"))

services = %w(
  AmazonEC2
  AmazonElastiCache
  
