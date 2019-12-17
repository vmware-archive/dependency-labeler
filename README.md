# deplab

deplab generates and shows metadata about a container image's dependencies.

## Obtain

Download the latest deplab release matching your OS from https://github.com/pivotal/deplab/releases/latest. Make it executable and move it to a directory in your `PATH` renaming it `deplab`.

## Usage

By default `deplab` [generates](#generate-metadata) the metadata of an image and the provided git repository (from where the image is built). The metadata is placed in a label on the output image, which are read by Pivotal's Open Source Licensing (OSL) process when using the container_image scan root. Once an image is labelled with `deplab` the metadata can be visualized using [inspect](#inspect).

To generate the metadata and output a labelled image, run
```bash
./deplab --image <image-name> --git <path to git repo> --output-tar <path to output tar>
```

Then, to visualise the metadata (optional), run 
```bash
./deplab inspect --image-tar <path to output tar>
```

## Generate metadata

`deplab` requires two input flags: an image source (remote `--image` or a local archive `--image-tar`) and the `--git` flag. At least one output flag needs to be specified (`--output-tar`, `--metadata-file`, `--dpkg-file`).  

```bash
./deplab --image-tar <path to input tar> \
    --git <path to git repo> \
    --output-tar <path to output tar>
```

### Generate flags

| short flag  | long flag  | value type | description | remarks |
|---|---|---|---|---|
| `-g` | `--git` | path |  [path to a directory under git revision control](#git) | Required. Can be provided multiple times. | 
| `-i` | `--image` | string | [image which will be analysed by deplab](#image) | Optional. Cannot be used with `--image-tar` flag | 
| `-p` | `--image-tar` |  path | [path to tarball of input image](#image-tarball) | Optional, but required for Concourse. Cannot be used with `--image` flag | 
| `-u` | `--additional-source-url` | url |  [url to the source of a dependency](#additional-source-url) | Optional. Can be provided multiple times. | 
| `-a` | `--additional-sources-file` | path |  [path to file containing yaml describing additional sources](#additional-sources-file) | Optional. Can be provided multiple times. | 
| `-t` | `--tag` | string | [tags the output image](#tag) | Optional | 
| `-d` | `--dpkg-file` | path | [write dpkg list metadata in (modified) '`dpkg -l`' format to a file at this path](#dpkg-file)| Optional |
| `-m` | `--metadata-file` | path | [write metadata to this file at the given path](#metadata-file) | Optional | 
| `-o` | `--output-tar` | path | [path to write a tarball of the image to](#tar) | Optional, but required for Concourse | 
|  | `--ignore-validation-errors` |  | By default deplab will exit with a non-zero exit code if a validation error is encountered. This flag will instead force deplab to output the validation failure message as a warning in StdErr and continue.  | Optional | 
| `-h` | `--help` |  | help for deplab |  | 
|  | `--version` |  |  version for deplab |  | 

## Inspect
Inspect is used to view the deplab metadata on a container image. Inspect prints the deplab "io.pivotal.metadata" label in the config file of an image to stdout.  The label will be printed in JSON format. If metadata does not exist on the image an error will be printed to standard error.

`deplab inspect` requires one image source to be specified (`--image` or `--image-tar`).

```bash
./deplab inspect --image <image-name>
```

### Inspect flags

| short flag  | long flag  | value type | description | remarks |
|---|---|---|---|---|
| `-i` | `--image` | string | [image to be inspected by deplab](#image) | Optional. Cannot be used with `--image-tar` flag | 
| `-p` | `--image-tar` |  path | [path to tarball of input image to be inspected by deplab](#image-tarball) | Optional, but required for Concourse. Cannot be used with `--image` flag | 

## Detailed flag descriptions

### Input flag descriptions

#### Git

You can specify as many git repositories as required by passing more than one
git flag into the command.

#### Image

deplab accepts as input an image stored in the local registry (tags, sha, or image id are all valid options).
One and only one of `--image` or `--image-tar` have to be used when invoking deplab.

#### Image tarball

deplab accepts as input an image stored in tar format (e.g. the output of `docker save ...` or of a concourse task).
One and only one of `--image` or `--image-tar` have to be used when invoking deplab.

#### Additional sources
Your image may have additional dependencies installed. These are dependencies which cannot be interpreted by dpkg or have been specified using the `--git` flag.
For OSL purposes you need to provide the source of these dependencies. The flags below allow you to specify the sources for these dependencies.

##### Additional source url

Additional source url allows you to specify a url which points to an archived source of a dependency. You can specify as many source urls as required using additional `--additional-source-url` flags.

Validation: The urls must be valid and reachable.  There is also a check to ensure that the url points to a compressed file type. Only the extension is checked and not the contents of the file.  On encountering an invalid url, deplab will provide an error message in StdErr.  By default deplab will exit with a non-zero exit code.  This default behaviour can be altered by using the `--ignore-validation-errors` flag, and deplab will continue and exit with a zero exit code.

##### Additional sources file

Additional sources file allows you to specify sources for additional dependencies as source archives or version control systems. You can specify as many of each type as required within a file, and as many additional sources files as required by passing more than one `--additional-sources-file` flags.

Validation: 
* archives: The urls must be valid and reachable.  There is also a check to ensure that the url points to a compressed file type. Only the extension is checked and not the contents of the file.
 * vcs: The url for git repository urls must start with on of the following: git:, ssh:, http:, https: or git@xxxx. 
 On encountering an invalid value, deplab will provide an error message in StdErr.  By default deplab will exit with a non-zero exit code.  This default behaviour can be altered by using the `--ignore-validation-errors` flag, and deplab will continue and exit with a zero exit code.

Supported format of the yaml file:
```yaml
archives:
- url: <url to source archive>
- url: <url to source archive>
vcs:
- protocol: git
  commit: <commit sha>
  url: <git repository url>
```

### Output flag descriptions

#### Tag

Optionally, the image can be tagged when exported as tar using the provided tag. The tag needs to be a valid docker tag.

#### Tar

Optionally deplab can output the image in tar format.

If a file exists at the given path, the file will be overwritten.

#### Metadata file

Optionally deplab can output the metadata to a file providing the path with the argument `--metadata-file` or `-m` 

If a file exists at the given path, the file will be overwritten.

#### dpkg file

Optionally deplab can output the debian package list portion of the metadata to a file with the argument `--dpkg-file` or `-d`

If a file exists at the given path, the file will be overwritten.

This file is approximately similar to the file which will be output by running `dpkg -l`, with the addition of an extra header which provides an ID for this list.

## Examples

### Basic usage

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --output-tar <path-to-image-output> 
```

### Multiple git inputs

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --git <path-to-another-repo> \
  --output-tar <path-to-image-output> 
```

### Input image as tar

```
deplab --image-tar <path-to-image-tar> \
  --git <path-to-repo> \
  --output-tar <path-to-image-output> 
```

### Multiple additional source url inputs

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --additional-source-url <url to archive> \
  --additional-source-url <url to archive> \
  --output-tar <path-to-image-output>
```

### Multiple additional-sources-file inputs

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --additional-sources-file <path to file> \
  --additional-sources-file <path to file> \
  --output-tar <path-to-image-output>
```

### Tag output image

```
deplab --image <image-reference> \
  --git <path-to-repo> \
  --tag <tag> \
  --output-tar <path-to-image-output>
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

### inspecting a tarball file

```
deplab inspect --image-tar <image-reference>
```

## Usage in Concourse

Please see [CONCOURSE.md](CONCOURSE.md) for information about using deplab as a task in your
Concourse pipeline.
 
## Data

##### debian package list

The `debian_package_list` requires the Debian package db to be present at `/var/lib/dpkg/status` or `/var/lib/status.d/*` on the image being instrumented on.
If not present, the dependency of type `debian_package_list` will be omitted.

`version` contains the _sha256_ of the `json` content of the metadata. Successive run of deplab on containers with the same `packages` and `apt_sources` will generate the same digest.

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

##### additional source url

For each `--additional-source-url` flag provided an archive object will be present in the metadata

```json
{
  "dependencies": [
    {...},
    {
      "type": "package",
      "source": {
        "type": "archive",
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
Testing requires `go` to be installed. To run tests you need to be authenticated against `dev.registry.pivotal.io` and be authorized for read access to `navcon/deplab-test-asset` repository.
```bash
go test ./...
```

Tests that pull images from registry are tagged `[remote-image]`. Tests that pull from a private registry that require authentication are tagged `[private-registry]`.

To skip tests you can run 

```bash
go test ./...  -ginkgo.skip='\[private-registry\]'
```

## Building

To build for release, please run the following:
```bash
go build -o deplab ./cmd/deplab
```

## Support

This tool is currently maintained by the Pivotal NavCon team; #navcon-team channel on Pivotal Slack.

Please reach out to us on Slack first, and then raise a Github issue.
