#!/bin/bash

if ! command -v yq &> /dev/null; then
    echo "Error: yq is not installed. Install it using:"
    echo "wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq"
    exit 1
fi

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

PROXY_FILE=$(realpath "$SCRIPTPATH/../out/all_proxies.yaml")
TEST_URL="https://ifconfig.me/ip"
TIMEOUT=10

if [ ! -f "$PROXY_FILE" ]; then
    echo "Error: File $PROXY_FILE not found"
    exit 1
fi

check_proxy() {
    local ip=$1
    local port=$2
    local protocol=$3
    
    echo "Checking $protocol://$ip:$port"
    
    if curl --proxy "$protocol://$ip:$port" \
            --max-time $TIMEOUT \
            --silent \
            --output /dev/null \
            --write-out "%{http_code}" \
            "$TEST_URL" | grep -q "200"; then
        echo "✅ Proxy working"
        return 0
    else
        echo "❌ Proxy not working"
        return 1
    fi
}

PROXY_COUNT=$(yq '. | length' "$PROXY_FILE")

for ((i=0; i<$PROXY_COUNT; i++)); do
    IP=$(yq ".[$i].ip" "$PROXY_FILE")
    PORT=$(yq ".[$i].port" "$PROXY_FILE")
    PROTOCOL=$(yq ".[$i].protocol" "$PROXY_FILE")
    
    echo "-------------------"
    echo "Testing proxy #$((i+1))"
    check_proxy "$IP" "$PORT" "$PROTOCOL"
done

echo "-------------------"
echo "Testing completed"
