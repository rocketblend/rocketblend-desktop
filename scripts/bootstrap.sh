#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

docker build -t rocketblend_desktop_builder .
docker run --rm -v $(pwd)/build/bin:/artifacts rocketblend_desktop_builder