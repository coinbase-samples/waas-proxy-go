
# Note: do not pass in pools/ID - just ID

POOL_ID=$1

URL="http://localhost:8443/v1/waas/proxy/pools/$POOL_ID"

curl -H "Content-Type: application/json" -X GET $URL
