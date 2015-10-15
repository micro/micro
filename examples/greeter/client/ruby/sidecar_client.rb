require 'net/http'
require 'json'

# Sidecar Client example in ruby
#
# This speaks to the service go.micro.srv.greeter
# via the sidecar application HTTP interface

uri = URI("http://localhost:8081/rpc")
service = "go.micro.srv.greeter"
method = "Say.Hello"
request = {"name" => "John"}

# do request
rsp = Net::HTTP.post_form(uri, {
	"service" => service,
	"method" => method,
	"request" => request
})

puts JSON.parse(rsp.body)["msg"]
