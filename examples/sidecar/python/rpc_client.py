import requests
import json
import sidecar

def main():
    response = sidecar.rpc_call("go.micro.srv.greeter", "Say.Hello", {"name": "John"})
    print response

if __name__ == "__main__":
    main()
