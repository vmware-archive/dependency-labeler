FROM golang:1.12.5 as builder

WORKDIR /workspace
# Copy the go source
COPY . /workspace

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o out cmd/deplab/main.go

FROM pivotalnavcon/docker-go

COPY --from=builder /workspace/out /deplab
ADD task/run.sh /usr/local/bin/deplab
