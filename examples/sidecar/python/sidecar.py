import uuid
import requests
import json

registry_uri = "http://localhost:8081/registry"
rpc_uri = "http://localhost:8081/rpc"
http_uri = "http://localhost:8081"
headers = {'content-type': 'application/json'}

def register(service):
    return requests.post(registry_uri, data=json.dumps(service), headers=headers)

def deregister(service):
    return requests.delete(registry_uri, data=json.dumps(service), headers=headers)

def rpc_call(service, method, request):
    payload = {
	"service": service,
        "method": method,
        "request": request,
    }
    return requests.post(rpc_uri, data=json.dumps(payload), headers=headers).json()

def http_call(path, request):
    return requests.post(http_uri + path, data=request)

