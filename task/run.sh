#!/usr/bin/env bash

set -eu -o pipefail

source /opt/resource/common.sh

check_env_var_exists() {
    if [[ -z "${2}" ]]; then
        echo "You must set the param ${1}"
        exit 1
    fi
}

main() {
    check_env_var_exists "IMAGE_TAR" "${IMAGE_TAR}"
    check_env_var_exists "OUTPUT_DIR" "${OUTPUT_DIR}"
    check_env_var_exists "GIT_REPO" "${GIT_REPO}"

    start_docker 3 3 "" ""

    /deplab --image-tar ${IMAGE_TAR} \
     --git ${GIT_REPO} \
     --output-tar ${OUTPUT_DIR}/image.tar
}

main