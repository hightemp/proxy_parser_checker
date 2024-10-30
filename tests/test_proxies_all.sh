#!/bin/bash

curl -X GET \
     -H "Content-Type: application/json" \
     http://localhost:8081/proxies/all
echo