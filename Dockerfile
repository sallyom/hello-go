# BUILD STAGE
FROM registry.access.redhat.com/ubi9/go-toolset as builder
USER root
ENV GOPATH=/opt/app-root GOCACHE=/mnt/cache GO111MODULE=on
WORKDIR $GOPATH/src/github.com/golang-ex
COPY . .
RUN go build -o hello-go ./hello.go

# RUN STAGE
FROM registry.access.redhat.com/ubi9/ubi-micro
ARG ARCH=amd64
COPY --from=builder /opt/app-root/src/github.com/golang-ex/hello-go /usr/bin/hello-go
ENTRYPOINT ["/usr/bin/hello-go"]
