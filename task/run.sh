#!/usr/bin/env bash

source /opt/resource/common.sh

start_docker 3 3 "" ""

/deplab --image-tar ${IMAGE_TAR} \
 --git ${GIT_REPO} \
 --output-tar ${OUTPUT_DIR}/image.tar