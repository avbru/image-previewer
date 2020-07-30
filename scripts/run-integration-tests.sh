#!/bin/bash

set -e

cd deployments

trap "docker-compose -p integration-test -f docker-compose.tests.yaml down" EXIT
docker-compose -p integration-tests -f docker-compose.tests.yaml run integration-test