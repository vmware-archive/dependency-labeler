#!/usr/bin/env bash

set -eu -o pipefail
set -x

for dockerfile in integration/assets/Dockerfile* ; do
    (cd integration/assets
      filename=$(basename $dockerfile)
      image_name=pivotalnavcon/ubuntu-${filename/Dockerfile./}

      docker build . -f $filename -t ${image_name}
      docker push ${image_name}
    )
done
