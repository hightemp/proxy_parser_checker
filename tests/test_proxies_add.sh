#!/bin/bash

curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
       "ip": "192.168.1.1",
       "port": "8080",
       "protocol": "http"
     }' \
     http://localhost:8080/proxies/add
echo