# `deplab` task

A Concourse task for adding metadata to Containers for OSL.

<!-- toc -->

- [usage](#usage)
  * [`params`](#params)
  * [`inputs`](#inputs)
  * [`outputs`](#outputs)
  * [`run`](#run)
- [example](#example)

<!-- tocstop -->

## usage

The task implementation is available as an image on Docker Hub at
[`pivotalnavcon/deplab-task`](http://hub.docker.com/r/pivotalnavcon/deplab-task). (This
image is built from [`Dockerfile.task`](Dockerfile.task).)

### `params`

Next, all of the following required parameters must be specified:

* `$IMAGE_TAR`: the path to the image to be labeled. The image must be in tarball format.

* `$OUTPUT_DIR`: the path to write the image to.

* `$GIT_REPOS`: The path to the git repo from which the git metadata will be generated. This should be
the source code of the application in the labelled image, and passed from an image build step to ensure
the correct commit SHA is provided. This can be a space separated list to allow for multiple git repositories.
ß
### `inputs`

There are two required inputs - a source for the image tarball and a space separated list of git repos:

```yaml
params:
  IMAGE_TAR: image/image.tar
  GIT_REPOS: git-deplab/

inputs:
- name: image
- name: git-deplab
```

### `outputs`

A single output may be configured:

```yaml
params:
  OUTPUT_DIR: labelled-image

outputs:
- name: labelled-image
```

The output will contain the following files:

* `image.tar`: the OCI image tarball. This tarball can be uploaded to a
  registry using the [Registry Image
  resource](https://github.com/concourse/registry-image-resource#out-push-an-image-up-to-the-registry-under-the-given-tags).
* `image-metadata.json`: the metadata in json format which has already been added to the OCI image tarball as a label.
* `image-dpkg-list.txt`: the debian package list portion of the metadata in `dpkg -l` format with additional headers

### `run`

Your task should run the `build` executable:

```yaml
run:
  path: deplab
```

## example

```yaml
  - task: deplab
    privileged: true
    config:
      platform: linux

      image_resource:
        type: docker-image
        source:
          repository: pivotalnavcon/deplab-task

      params:
        IMAGE_TAR: image/image.tar
        GIT_REPOS: git-deplab
        OUTPUT_DIR: labelled-image

      inputs:
      - name: git-deplab
      - name: image

      outputs:
      - name: labelled-image

      run:
        path: deplab
```