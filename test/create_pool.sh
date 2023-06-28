
URL="http://localhost:8443/v1/waas/proxy/pools"

POOL_ID=$(uuidgen)

curl -H "Content-Type: application/json" -X PUT -d "{ \"pool\": { \"display_name\": \"$POOL_ID\"} }" $URL
