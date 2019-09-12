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
    check_env_var_exists "GIT_REPOS" "${GIT_REPOS}"

    if [[ -z "${DPKG_FILE}" ]]; then
      DPKG_FILE="image-dpkg-list.txt"
    fi

    start_docker 3 3 "" ""

    args=""
    for repo in ${GIT_REPOS}; do
      args+="--git ${repo} "
    done

    /deplab --image-tar ${IMAGE_TAR} \
     --output-tar ${OUTPUT_DIR}/image.tar \
     --metadata-file ${OUTPUT_DIR}/image-metadata.json \
     --dpkg-file ${OUTPUT_DIR}/${DPKG_FILE} \
     ${args[@]}
}

main