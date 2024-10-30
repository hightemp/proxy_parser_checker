#!/bin/bash

curl -X GET \
     -H "Content-Type: application/json" \
     http://localhost:8080/proxies/all
echo