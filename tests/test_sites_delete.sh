#!/bin/bash

curl -X DELETE \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://example.com/proxy-list"
     }' \
     http://localhost:8081/sites/delete
echo 