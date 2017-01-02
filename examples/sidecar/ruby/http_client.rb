require './sidecar'

puts http_call("/greeter", {"name" => "John"})
