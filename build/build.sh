#!/usr/bin/env bash
DIR=$(cd `dirname $0` && pwd)
PARENTDIR="$(dirname "${DIR}")"
DOCKERFILE=${DIR}/package/Dockerfile
TAG=$(git describe --abbrev=0 --tags)
docker build -f ${DOCKERFILE} -t milightd:${TAG:1} ${PARENTDIR}
