# deplab

## Introduction
deplab adds metadata about a container image's dependencies as a label to the container image.

## Dependencies
Docker is required to be installed and available on your path, test this by running `docker version`.
API version 1.39 or higher is required.

## Usage
Download the latest `deplab` binary from the [releases page](https://github.com/pivotal/deplab/releases).
To run the tool run the following command:
```bash
./deplab [flags]
```
This returns the sha256 or the provided tag of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

To visualise the metadata this command can be run

```bash
docker inspect $(./deplab --image <image-name> --git <path to git repo>) \
  | jq -r '.[0].Config.Labels."io.pivotal.metadata"' \ 
  | jq .
```

## Flags

| short flag  | long flag  | value type | description | remarks |
|---|---|---|---|---|
| `-g` | `--git` | path |  [path to a directory under git revision control](#git) | Required. Can be provided multiple times. | 
| `-i` | `--image` | string | [image which will be analysed by deplab](#image) | Optional. Cannot be used with `--image-tar` flag | 
| `-p` | `--image-tar` |  path | [path to tarball of input image](#image-tarball) | Optional, but required for Concourse. Cannot be used with `--image` flag | 
| `-b` | `--blob` | url |  [url to the source of a dependency](#blob) | Optional. Can be provided multiple times. | 
| `-t` | `--tag` | string | [tags the output image](#tag) | Optional | 
| `-d` | `--dpkg-file` | path | [write dpkg list metadata in (modified) '`dpkg -l`' format to a file at this path](#dpkg-file)| Optional |
| `-m` | `--metadata-file` | path | [write metadata to this file at the given path](#metadata-file) | Optional | 
| `-o` | `--output-tar` | path | [path to write a tarball of the image to](#tar) | Optional, but required for Concourse | 
| `-h` | `--help` |  | help for deplab |  | 
|  | `--version` |  |  version for deplab |  | 

### Inputs

#### Git

You can specify as many git repositories as required by passing more than one
git flag into the command.

#### Image

deplab accept as input an image stored in the local registry (tags, sha, or image id are all valid options).
One and only one of `--image` or `--image-tar` have to be used when invoking deplab.

#### Image tarball

deplab accept as input an image stored in tar format (e.g. the output of `docker save ...` or of a concourse task).
One and only one of `--image` or `--image-tar` have to be used when invoking deplab.

#### Blob

Blob is to allow any arbitrary url which points to a dependency source. You can specify as many blob urls as required by passing more than one
blob flag into the command.

### Outputs

#### Default

This returns the sha256 of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

#### Tag

Optionally image can be tagged in the local registry using the provided tag. Tag need to be a valid docker tag.

#### Tar

Optionally deplab can output the image in tar format.

If the file path cannot be created deplab will process the image and store it in Docker, but will also return an error for the writing of the tar. 

If a file exists at the given path, the file will be overwritten.

#### Metadata file

Optionally deplab can output the metadata to a file providing the path with the argument `--metadata-file` or `-m` 

If the file path cannot be created, deplab will return the newly labelled image, and return an error for the writing of the metadata file. 

If a file exists at the given path, the file will be overwritten.

#### dpkg file

Optionally deplab can output the debian package list portion of the metadata to a file with the argument `--dpkg-file` or `-d`

If the file path cannot be created, deplab will return the newly labelled image, and return an error for the writing of the dpkg file. 

If a file exists at the given path, the file will be overwritten.

This file is approximately similar to the file which will be output by running `dpkg -l`, with the addition of an extra header which provides an ID for this list.


## Examples

### Basic usage

```
deplab --image <image-reference> \
  --git <path-to-repo> 
```

### Multiple git inputs

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --git <path-to-another-repo>
```

### Input image as tar

```
deplab --image-tar <path-to-image-tar> \
  --git <path-to-repo> 
```

### Output image as tar

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --output-tar <path-to-image-output> 
```


### Multiple blob inputs

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --blob <url to blob> \
  --blob <url to blob>
```

### Tag output image

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --tag <tag> 
```

### dpkg list file

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --dpkg-file <path-to-dpkg-file-output> 
```

### metadata file

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --metadata-file <path-to-metadata-file-output> 
```

### Usage in Concourse

Please see [CONCOURSE.md](CONCOURSE.md) for information about using deplab as a task in your
Concourse pipeline.
 
## Data

##### debian package list

The `debian_package_list` requires the Debian package db to be present at `/var/lib/dpkg/status` or `/var/lib/status.d/*` on the image being instrumented on.
If not present, the dependency of type `debian_package_list` will be omitted.

`version` contains the _sha256_ of the `json` content of the metadata. Successive run of deplab on containers with the same `packages` and `apt_sources` are going to generate the same digest.

The debian package list is generated with the following format.

```json
{
  "dependencies": [
    {
      "type": "debian_package_list",
      "version": {
        "sha256": "a56...42b"
      },
      "source": {
        "type": "inline",
        "version": null,
        "metadata": {
          "packages": [...],
          "apt_sources": [...]
        }
      }
    }
  ]
}
```

Example of a package item in field `packages` 

```json
{
  "package": "zlib1g",
  "version": "1:1.2.11.dfsg-0ubuntu2",
  "architecture": "amd64",
  "source": {
    "package": "zlib",
    "version": "1:1.2.11.dfsg-0ubuntu2",
    "upstreamVersion": "1.2.11.dfsg"
  }
}
```

Example of `apt_sources` content

```json
[
  "deb http://archive.ubuntu.com/ubuntu/ bionic main restricted",
  "deb http://archive.ubuntu.com/ubuntu/ bionic-updates main restricted",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security main restricted",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security universe",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security multiverse"
]
```

##### git dependency
   
   For each `--git` flag provided a git dependency will be present in the metadata
   
   If the `--git` flag is provided with a valid path to a git repository, a git dependency will be added:
   ```json
   {
     "dependencies": [
       {...},
       {
         "type": "package",
         "source": {
           "type": "git",
           "version": {
             "commit":  "d2c[...]efd"
            },
           "metadata": {
             "url": "https://github.com/pivotal/deplab.git",
             "refs": ["0.5.0"]
           }
         }
       }
     ]
   }
   ```

##### blob

For each `--blob` flag provided a blob object will be present in the metadata

```json
{
  "dependencies": [
    {...},
    {
      "type": "package",
      "source": {
        "type": "blob",
        "metadata": {
          "url": "http://archive.ubuntu.com/ubuntu/pool/main/c/ca-certificates/ca-certificates_20180409.tar.xz"
        }
      }
    }
  ]
}
```

#### base
The base image metadata is generated with the following format
```json
  "base": {
    "name": "Ubuntu",
    "version_id": "18.04",
    "version_codename": "bionic",
    ...
  }
```

it includes all the content of `/etc/os-release` present on the image (keys are lower-cased).

This relies on the `/etc/os-release` file being in the docker container. If `/etc/os-release` is not present all the field will be set to `unknown`.

```json
{
  "name": "unknown",
  "version_id": "unknown",
  "version_codename": "unknown"
}
```

#### provenance
Provenance is a list of the tools which have added information to the image. It is generated in the following format
```json
  "provenance": [
    {
      "name": "deplab",
      "version": "0.0.0-dev",
      "url": "https://github.com/pivotal/deplab"
    }
  ]
```


## Testing
Testing requires `go` to be installed.  Please clone this git repository.  Tests can be run with:
```bash
go test ./...
```

## Building

To build for release, please run the following:
```bash
go build -o deplab ./cmd/deplab
```

To build the Concourse task image, please run the following:
```bash
docker build . -f Dockerfile.task
```

## Support

This tool is currently maintained by the Pivotal NavCon team;
@navcon in #navcon-team on Pivotal Slack.

Please reach out to us on Slack first, and then raise a Github issue.
