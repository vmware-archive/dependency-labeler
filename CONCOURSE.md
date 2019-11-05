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

The task implementation is available as an image on Pivotal's internal registry at
[`dev.registry.pivotal.io/navcon/deplab-task`](https://dev.registry.pivotal.io/harbor/projects/15/repositories). (This
image is built from [`Dockerfile.task`](Dockerfile.task).)

To gain access to this image, you will have to request authorization for you (or your bot account).

1. Please ensure you (or your bot) has a PivNet account
1. Send a slack message to @navcon in [#navcon-team](https://app.slack.com/client/T024LQKAS/CFUA5BXV5/thread/C2Y1X7ZAN-1572563183.022900) requesting access to the `deplab-task` image, providing the email address corresponding to your (or your bot's) PivNet account.
1. Once NavCon have confirmed that you (or your bot) have been added to the `deplab-task` repository, you can access the image using your (or your bot's) PivNet username and password.

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
          repository: dev.registry.pivotal.io/navcon/deplab-task
          username: ((your-pivnet-username))
          password: ((your-pivnet-password))

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