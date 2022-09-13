#!/bin/bash

export warpgate_data="/tmp/warpgate/data" # randomize the name 

docker compose -f "$(pwd)"/_scripts/docker-compose.yml down

sudo rm -r ${warpgate_data}