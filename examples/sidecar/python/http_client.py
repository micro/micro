import requests
import json
import sidecar

def main():
    response = sidecar.http_call("/greeter", {"name": "John"})
    print response.text

if __name__ == "__main__":
    main()
