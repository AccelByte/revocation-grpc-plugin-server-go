#!/usr/bin/env bash

BIN_SERVER_NAME=revocation_server
BIN_DEMO_NAME=revocation_demo_cli

echo "[*] Build sample app"
go build -o $BIN_SERVER_NAME .
(cd demo/cli && go build -o ../../$BIN_DEMO_NAME .)

./$BIN_SERVER_NAME &
(for _ in $(seq 1 10); do bash -c "timeout 1 echo > /dev/tcp/127.0.0.1/8080" 2>/dev/null && exit 0 || sleep 3s; done; exit 1)

echo "[*] Sample app running"

echo "[*] Run ngrok tcp 6565"
( ngrok tcp 6565 > /dev/null 2>&1 & )

for _ in {1..10}
do
    RESP=$(curl -s --location 'localhost:4040/api/tunnels')
    SERVER_URL=$(echo "$RESP" | jq -r '.tunnels[] | select(.config.addr = "localhost:6565") | .public_url')
    if [ -z "$SERVER_URL" ]; then
        echo "[?] ngrok still waiting..."
        sleep 5s
    else
        echo "[*] ngrok online on: $SERVER_URL"
        break
    fi
done

export GRPC_SERVER_URL="${SERVER_URL#*://}"
echo "[*] Run demo/cli test"
./$BIN_DEMO_NAME

rm $BIN_DEMO_NAME $BIN_SERVER_NAME