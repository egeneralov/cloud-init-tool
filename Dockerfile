FROM golang
ENV GO111MODULE=on
WORKDIR /go/src/github.com/egeneralov/cloud-init-tool
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build -o /go/bin/cloud-init

FROM alpine
ENV PATH=$PATH:/go/bin/
COPY --from=0 /go/bin/cloud-init /go/bin/cloud-init
