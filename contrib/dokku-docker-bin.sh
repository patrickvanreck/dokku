#!/usr/bin/env bash
set -eo pipefail
[[ $DOKKU_TRACE ]] && set -x && DOKKU_DOCKER_ENV="--env DOKKU_TRACE=on"
DOKKU_DOCKER_CONTAINER_NAME=${DOKKU_DOCKER_CONTAINER_NAME:=dokku}

# TODO: handle cases where we need a tty
# shellcheck disable=SC2086
docker exec $DOKKU_DOCKER_ENV -i "$DOKKU_DOCKER_CONTAINER_NAME" dokku "$@"
