#!/usr/bin/env bash

set -eu -o pipefail
set -x

for dockerfile in $(ls integration/assets) ; do
    image_name=pivotalnavcon/ubuntu-${dockerfile/Dockerfile./}

    docker build . -f integration/assets/$dockerfile -t ${image_name}
    docker push ${image_name}
done
