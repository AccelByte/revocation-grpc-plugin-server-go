#!/usr/bin/env bash

# this script should be run inside docker container with
# prerequisites: bash, go, ngrok, curl, jq

BIN_SERVER_NAME=sample_app_server
BIN_DEMO_NAME=sample_app_demo_cli

echo "- [*] Build sample app"
go build -o $BIN_SERVER_NAME .
(cd demo/cli && go build -o ../../$BIN_DEMO_NAME .)

./$BIN_SERVER_NAME &
(for _ in {1..10}; do bash -c "timeout 1 echo > /dev/tcp/127.0.0.1/8080" 2>/dev/null && exit 0 || sleep 3s; done; exit 1)
if [ $? -ne 0 ]
then
    echo "- [!] Failed to run sample app"
    exit 1
fi

echo "- [*] Sample app running"

echo "- [*] Run ngrok tcp 6565"
( ngrok tcp 6565 > tmp.dat 2>&1 & )

for _ in {1..1}
do
    sleep 3s
    RESP=$(curl -s --location 'localhost:4040/api/tunnels')
    SERVER_URL=$(echo "$RESP" | jq -r '.tunnels[] | select(.config.addr = "localhost:6565") | .public_url')
    if [ -n "$SERVER_URL" ]; then
        echo "- [*] ngrok online on: $SERVER_URL"
        break
    fi
done

[ -z "$SERVER_URL" ] && echo "- [!] Failed to run ngrok" && cat tmp.dat && rm tmp.dat && exit 1

export GRPC_SERVER_URL="${SERVER_URL#*://}"
echo "- [*] Run demo/cli test"
./$BIN_DEMO_NAME

rm $BIN_DEMO_NAME $BIN_SERVER_NAME tmp.dat