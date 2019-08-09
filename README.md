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
./deplab --image <image name>
```

Where `<image name>` is the name of the image that you want to add the meta data to.

This returns the sha256 of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

## Data
The deplab tool adds the following dependencies:
- debian_package_list

##### debian_package_list
The debian_package_list is generated with the following format

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
