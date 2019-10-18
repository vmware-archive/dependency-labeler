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

There are no params, instead the run.args should be used to pass in deplab flags

### `inputs`

There are two required inputs - a source for the image tarball (referenced in the `--image-tar` flag) and one (or more) git repositories from where the image was built (referenced in the `--git` flags):

```yaml
inputs:
- name: image
- name: git-deplab
```

### `outputs`

An output should be configured for the deplab files which need to be saved (i.e. the `--output-tar`, `--metadata-file`, and the `--dpkg-file` flags):

```yaml
outputs:
- name: labelled-image
```

### `run`

Your task should run the `deplab` executable, and provide the [command line args expected by deplab](README.md):

```yaml
run:
  path: deplab
  args:
  - --image-tar
  - image/image.tar
  - --git
  - git-deplab
  - --output-tar
  - labelled-image/image.tar
```

## example

```yaml
  - task: deplab
    config:
      platform: linux

      image_resource:
        type: docker-image
        source:
          repository: pivotalnavcon/deplab-task

      inputs:
      - name: git-deplab
      - name: image

      outputs:
      - name: labelled-image

      run:
        path: deplab
        args:
        - --image-tar
        - image/image.tar
        - --git
        - git-deplab
        - --output-tar
        - labelled-image/image.tar
```


## example using interpolated values
 
Sometime you may need to interpolate some of your values.  In the example below we are using the version as a 
customisable run time variable.

```yaml
run:
  path: /bin/bash
  args:
    - -cex
    - |
      version="$(cat version-navcon-test-app/version)"
      deplab \
        --image-tar image/image.tar \
        --git git-navcon-test-app \
        --output-tar "annotated-image/image-$version.tar" \
        --dpkg-file annotated-image/navcon-test-app-dpkg-list.txt \
        --tag "$version"
```