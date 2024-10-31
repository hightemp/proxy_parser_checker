#!/bin/bash

curl -X DELETE \
     -H "Content-Type: application/json" \
     -d '{
       "ip": "192.168.1.1",
       "port": "8080",
       "protocol": "http"
     }' \
     http://localhost:8081/proxies/delete
echo 