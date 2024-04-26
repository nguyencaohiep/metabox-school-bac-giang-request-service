#!/bin/bash

# ./docker-build.sh tag_image

tag_image=$1
if [ -z $1 ]; then
    echo "Chua nhap tham so, gia tri se de mac dinh tag:latest"
    tag_image=latest
fi

# build go
go build -ldflags "-s -w" -trimpath -o gift ./cmd/*.go

# build image
docker build -f deployments/Dockerfile_build_local -t registry.ctisoftware.vn/outsource/gift:"$tag_image" .

# remove go-bin-file
rm gift

# shellcheck disable=SC2046
docker rmi -f $(docker images -f dangling=true -q)