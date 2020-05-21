module github.com/pivotal/deplab

go 1.12

require (
	github.com/Microsoft/hcsshim v0.8.6 // indirect
	github.com/containerd/containerd v1.3.1-0.20191007173812-901bcb223146
	github.com/containerd/continuity v0.0.0-20190827140505-75bee3e2ccb6 // indirect
	github.com/docker/docker v1.4.2-0.20180531152204-71cd53e4a197
	github.com/docker/go-units v0.4.0 // indirect
	github.com/google/go-containerregistry v0.0.0-20191008160043-1e84d6375038
	github.com/joho/godotenv v1.3.0

	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	golang.org/x/text v0.3.2
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.3.0

)

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
