#!/usr/bin/env bash

set -eu -o pipefail
set -x

for dockerfile in test/integration/assets/dockerfiles/Dockerfile* ; do
    (cd test/integration/assets/dockerfiles
      filename=$(basename $dockerfile)
      image_name=dev.registry.pivotal.io/navcon/deplab-test-asset:${filename/Dockerfile./}

      docker build . -f $filename -t ${image_name}
      docker push ${image_name}
    )
done
