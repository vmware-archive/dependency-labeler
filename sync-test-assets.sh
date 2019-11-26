#!/usr/bin/env bash

### requires to be logged in the target registry

set -eux -o pipefail

DEPLAB_ASSET_REPOSITORY="dev.registry.pivotal.io/navcon/deplab-test-asset"
IMAGE_ARCHIVES="test/integration/assets/image-archives"

### Build and push all dockerfiles to registry

for dockerfile in test/integration/assets/dockerfiles/Dockerfile* ; do
    (cd test/integration/assets/dockerfiles
      filename=$(basename "$dockerfile")
      image_name=$DEPLAB_ASSET_REPOSITORY:${filename/Dockerfile./}

      docker build . -f "$filename" -t "${image_name}"
      docker push "${image_name}"
    )
done

### Save locally image from registry

images=( "cloudfoundry/run:tiny" \
          "$DEPLAB_ASSET_REPOSITORY:all-file-types" \
          "$DEPLAB_ASSET_REPOSITORY:broken-files" \
          "$DEPLAB_ASSET_REPOSITORY:tiny-with-invalid-label" \
          "$DEPLAB_ASSET_REPOSITORY:os-release-on-scratch" )

for image in "${images[@]}"; do
  filename="${image##*:}.tgz"

  docker save "$image" -o "$IMAGE_ARCHIVES/$filename"
done

### Additional assets

go run cmd/deplab/main.go --image-tar "$IMAGE_ARCHIVES/tiny.tgz" --git . --output-tar "$IMAGE_ARCHIVES/tiny-deplabd.tgz"

crane push "$IMAGE_ARCHIVES/tiny-deplabd.tgz" "$DEPLAB_ASSET_REPOSITORY:tiny-deplabd"

tar czvf "$IMAGE_ARCHIVES/invalid-image-archive.tgz" "$IMAGE_ARCHIVES/../sources/empty-file.yml"