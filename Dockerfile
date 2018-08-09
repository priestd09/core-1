FROM golang:1.10.3-stretch AS build
RUN mkdir -p src/github.com/mesg-foundation/core
RUN go get github.com/xeipuuv/gojsonschema && \
    go get gopkg.in/yaml.v2 && \
    go get github.com/stvp/assert && \
    go get github.com/docker/docker/api/types/swarm && \
    go get github.com/docker/docker/api/types/mount && \
    go get github.com/docker/docker/client && \
    go get github.com/spf13/viper && \
    go get github.com/spf13/cobra && \
    go get gopkg.in/AlecAivazis/survey.v1 && \
    go get github.com/golang/protobuf/proto && \
    go get google.golang.org/grpc && \
    go get github.com/syndtr/goleveldb/leveldb && \
    go get github.com/cnf/structhash && \
    go get gopkg.in/src-d/go-git.v4/... && \
    go get github.com/asaskevich/govalidator && \
    go get github.com/cpuguy83/go-md2man && \
    go get github.com/mitchellh/go-homedir
ADD . src/github.com/mesg-foundation/core
WORKDIR src/github.com/mesg-foundation/core
RUN go get ./...
ARG version
RUN go build -o mesg-core \
      -ldflags="-X 'github.com/mesg-foundation/core/version.Version=$version'" \
      core/main.go

FROM ubuntu:18.04
WORKDIR /app
COPY --from=build /go/src/github.com/mesg-foundation/core/mesg-core .
CMD ["./mesg-core"]
