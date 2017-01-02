require './sidecar'

puts rpc_call("/greeter/say/hello", {"name": "John"})
