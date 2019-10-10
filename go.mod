module github.com/pivotal/deplab

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20180531152204-71cd53e4a197
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-containerregistry v0.0.0-20191008160043-1e84d6375038
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/jhoonb/archivex v0.0.0-20180718040744-0488e4ce1681
	github.com/joho/godotenv v1.3.0
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/nwaples/rardecode v1.0.0 // indirect

	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.7.0
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/build v0.0.0-20191008185658-817d966b7e93
	golang.org/x/text v0.3.2
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.4

)

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
