require 'open-uri'
require 'nokogiri'
require 'json'

url = 'https://aws.amazon.com/ec2/instance-types/'
data = {}

begin
  doc = Nokogiri::HTML(URI.open(url))

  doc.css('#aws-element-f137e899-0a8d-4358-a14d-f645975465a3 tr').each do |row|
    columns = row.css('td')
    next if columns.empty?
    next if columns.length <= 6
    instance_type = columns[0].text&.strip
    next if instance_type =~ /Instance/

    gpu_count = 0
    gpu_type = ""
    gpu_memory = 0
    gpu_memory_unit = ""

    case
    when instance_type.start_with?("p5")
      gpu_count = columns[2].text&.strip.split(" ").first.to_i
      gpu_type = "NVIDIA H100"
      gpu_memory = columns[4].text&.strip.split(" ", 2).first.to_i
      gpu_memory_unit = columns[4].text&.strip.split(" ", 2).last
    when instance_type.start_with?("p4")
      gpu_count = columns[2].text&.strip.to_i
      gpu_type = "NVIDIA A100"
      gpu_memory = columns[4].text&.strip.split(" ", 2).first.to_i
      gpu_memory_unit = columns[4].text&.strip.split(" ", 2).last

    when instance_type.start_with?("g6e")
      gpu_count = columns[3].text&.strip.to_i
      gpu_type = "NVIDIA L4"
      gpu_memory = columns[4].text&.strip.split(" ", 2).first.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("g5g.")
      gpu_count = columns[3].text&.strip.to_i
      gpu_type = "NVIDIA T4G"
      gpu_memory = columns[4].text&.strip.split(" ", 2).first.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("g5.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "NVIDIA A10G"
      gpu_memory = columns[2].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("g4dn.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "NVIDIA T4"
      gpu_memory = columns[4].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("g4ad.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "Radeon Pro V520"
      gpu_memory = columns[4].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("trn2.")
      gpu_count = columns[2].text&.strip.to_i
      gpu_type = "AWS Trainium2"
      gpu_memory = columns[5].text&.strip.to_i
      gpu_memory_unit = "TB"

    when instance_type.start_with?("trn1.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "AWS Trainium1"
      gpu_memory = columns[2].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("inf2.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "AWS Inferentia2"
      gpu_memory = columns[2].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("inf1.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "AWS Inferentia1"

    when instance_type.start_with?("dl1.")
      gpu_count = columns[2].text&.strip.to_i
      gpu_type = "Gaudi"

    when instance_type.start_with?("dl2q.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "Qualcomm AI 100"
      gpu_memory = columns[2].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("f2.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "FPGAs"
      gpu_memory = columns[3].text&.strip.to_i
      gpu_memory_unit = "GB"

    when instance_type.start_with?("vt1.")
      gpu_count = columns[1].text&.strip.to_i
      gpu_type = "U30"
      gpu_memory = columns[3].text&.strip.to_i
      gpu_memory_unit = "GB"
    end

    next if gpu_count == 0

    # Consider Trainium and Inferentia as GPU equivalents
    data[instance_type] = {
      'core' => gpu_count,
      'type' => gpu_type,
      'mem' => gpu_memory,
      'mem_unit' => gpu_memory_unit,
    }
  end

  puts JSON.pretty_generate(data)

rescue OpenURI::HTTPError => e
  puts "Error accessing the URL: #{e.message}"
rescue Nokogiri::CSS::SyntaxError => e
  puts "Error parsing the HTML: #{e.message}"
rescue StandardError => e
  puts "An unexpected error occurred: #{e.message}"
end