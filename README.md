# deplab

## Introduction
Deplab adds metadata about a container image's dependencies as a label to the container image.

## Dependencies
Docker is required to be installed and available on your path, test this by running `docker --version`.
Version 1.39 or higher is required.

## Usage
Download the latest `deplab` binary from the releases page.
To run the tool run the following command:
```bash
./deplab --image <image name> --git <path to git repo>
```

`<image name>` is the name of the image that you want to add the meta data to.

The `--git` flag is optional. The  `<path to git repo>` is a path to a directory under git version control.

This returns the sha256 of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

## Data

##### debian package list
The debian package list is generated with the following format

```json
{
  "dependencies": [
    {
      "type": "debian_package_list",
      "source": {
        "type": "inline",
        "version": null,
        "metadata": {
          "packages": [array_of_packages]
        }
      }
    }
  ]
}
```

The `debian_package_list` requires `dpkg` to be a package on the image being instrumented on. If not present, the metadata structure will be injected into metadata of the container, with an empty array of packages.

##### git dependency
If the `--git` flag is provided with a valid path to a git repository, a git dependency will be added:
```json
{
  "dependencies": [
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2c3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "uri": "https://github.com/pivotal/deplab.git",
          "refs": ["0.5.0"]
        }
      }
    }
  ]
}
```


## Testing
Testing requires `go` to be installed.  Please clone this git repository.  Tests can be run with:
```bash
go test
```

## Building

To build for release, please run the following:
```bash
go build -o deplab
```

## Support

This tool is currently maintained by the Pivotal NavCon team;
@navcon in #navcon-team on Pivotal Slack.

Please reach out to us on Slack first, and then raise a Github issue.
