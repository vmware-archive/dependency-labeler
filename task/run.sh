#!/usr/bin/env bash

source /opt/resource/common.sh

start_docker 3 3 "" ""

/deplab --image-tar ${IMAGE_TAR} \
 --git ${GIT_REPO} \
 --metadata-file ${OUTPUT_DIR}/metadata.json \
 --output-tar ${OUTPUT_DIR}/image.tar