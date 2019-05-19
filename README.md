# ES-PUSHER

[![Build Status](https://travis-ci.org/elastifeed/es-pusher.svg?branch=master)](https://travis-ci.org/elastifeed/es-pusher)
[![Docker Repository on Quay](https://quay.io/repository/elastifeed/es-pusher/status "Docker Repository on Quay")](https://quay.io/repository/elastifeed/es-pusher)

Provides a REST API and interfaces with Elasticsearch.

## Configuration
> Configuration is done via environment variables, see below for possible options
  - ES_ADDRESSES (e.g. `["http://host0:9200", "http://host1:9200", ...]`)
  - API_BIND (e.g. `:9000`)


## Running locally for debugging
`podman run --network=host -e ES_ADDRESSES='["http://localhost:9200"]' -e API_BIND=':9000' -d elastifeed/es-pusher:latest`

## Prometheus metrics
>  Metrics can be optained via `<ip:port>/metrics`
Available metrics:
- espusher_rest_calls
- espusher_rest_calls_malformed
- espusher_rest_calls_successful
- espusher_storage_elasticsearch_add_document_requested_count
- espusher_storage_elasticsearch_added_document_count
