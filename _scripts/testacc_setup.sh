#!/bin/bash

export warpgate_data="/tmp/warpgate/data" # TODO randomize the name 

sudo rm -rf ${warpgate_data}

mkdir -p ${warpgate_data}

function rand {
    openssl rand -base64 $@
}

export WARPGATE_HOST=127.0.0.1
export WARPGATE_PORT=38888
export WARPGATE_USERNAME=admin
export WARPGATE_PASSWORD=$(rand 16)
export WARPGATE_INSECURE_SKIP_VERIFY=true

docker run -it --rm -v ${warpgate_data}:/data ghcr.io/warp-tech/warpgate:v0.7.0 unattended-setup \
        --admin-password "${WARPGATE_PASSWORD}" \
        --data-path /data \
        --http-port "8888" 

docker-compose -f "$(pwd)"/_scripts/docker-compose.yml up -d --wait