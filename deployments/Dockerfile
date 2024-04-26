FROM golang:1.22-alpine as builder

LABEL stage=builder

ARG CGO_ENABLED=0
ARG GO111MODULE=on
ARG GIT_USERNAME=
ARG GIT_ACCESS_TOKEN=

ARG service_name=

WORKDIR /go/src/app

RUN apk add --no-cache openssh-client git
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN cd ./cmd && \
    go build -ldflags "-s -w" -trimpath -o app *.go

# SERVICE BUILDER
FROM alpine:latest

WORKDIR /opt/app-service

ENV ADDR_GRAPH 0.0.0.0:8081
ENV ADDR_GRPC 0.0.0.0:9091

COPY --from=builder /go/src/app/cmd/app .
COPY ./config/.env .

RUN apk update && apk add libc6-compat

RUN chmod +x ./app

ENV STUDENT_NUMBERRETRY 5
ENV STUDENT_EXPIRESTIMEACCESSTOKEN 1
ENV STUDENT_EXPIRESTIMEREFRESHTOKEN 3
ENV STUDENT_ACCESSTOKENTYPE access
ENV STUDENT_REFRESHTOKENTYPE refresh
ENV STUDENT_DATABASENAME student

CMD ./app --configSource=.env serve --addr-graph=$ADDR_GRAPH

