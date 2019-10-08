module github.com/pivotal/deplab

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.13 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/jhoonb/archivex v0.0.0-20180718040744-0488e4ce1681
	github.com/joho/godotenv v1.3.0
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect

	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.5.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.2
	gotest.tools v2.2.0+incompatible // indirect

)

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
