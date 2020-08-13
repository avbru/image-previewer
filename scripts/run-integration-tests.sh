#!/bin/bash

set -e

cd deployments

docker-compose -p integration-tests -f docker-compose.tests.yaml build --no-cache integration-test
docker-compose -p integration-tests -f docker-compose.tests.yaml run integration-test
docker-compose -p integration-tests -f docker-compose.tests.yaml down