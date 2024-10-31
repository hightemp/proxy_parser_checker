#!/bin/bash

curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://example.com/proxies"
     }' \
     http://localhost:8081/api/v1/sites
echo