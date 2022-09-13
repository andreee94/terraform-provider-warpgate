#!/bin/bash

export warpgate_data="/tmp/warpgate/data" # randomize the name 
warpgate_data_original="./_scripts/data"

# docker run --rm -it -v /tmp/warpgate/data:/data ghcr.io/warp-tech/warpgate setup

sudo rm -rf ${warpgate_data}

sudo mkdir -p ${warpgate_data}

sudo cp -r ${warpgate_data_original}/* ${warpgate_data}/

export WARPGATE_HOST=127.0.0.1
export WARPGATE_PORT=38888
export WARPGATE_USERNAME=admin
export WARPGATE_PASSWORD=password
export WARPGATE_INSECURE_SKIP_VERIFY=true

docker-compose -f "$(pwd)"/_scripts/docker-compose.yml up -d --wait