#!/bin/bash

curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://example.com/proxy-list"
     }' \
     http://localhost:8081/sites/add
echo