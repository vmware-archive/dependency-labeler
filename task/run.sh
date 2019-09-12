#!/usr/bin/env bash

set -eu -o pipefail

source /opt/resource/common.sh

main() {
    start_docker 3 3 "" ""

    /deplab $@
}

main $@