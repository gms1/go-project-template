#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

OPTIONS:
  -h|--help|help  ... display this usage information and exit
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage

echo "jaeger-start..."

docker rm -f jaeger &>/dev/null || true
docker volume prune -f &>/dev/null || true
docker run -d \
  -e COLLECTOR_OTLP_ENABLED=true \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 9411:9411 \
  --rm --name jaeger jaegertracing/all-in-one:1.64.0

# 16686	HTTP	query	serve frontend
#  4317	HTTP	collector	accept OpenTelemetry Protocol (OTLP) over gRPC
#  4318	HTTP	collector	accept OpenTelemetry Protocol (OTLP) over HTTP
# 14268	HTTP	collector	accept jaeger.thrift directly from clients
# 14250	HTTP	collector	accept model.proto
#  9411	HTTP	collector	Zipkin compatible endpoint (optional)

echo "jaeger-start: SUCCEEDED"
