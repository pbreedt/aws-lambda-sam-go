# see: https://gallery.ecr.aws/docker/library/golang for supported runtimes
FROM public.ecr.aws/docker/library/golang:1.21-alpine as build-image
WORKDIR /src
COPY go.mod go.sum main.go ./
# required for alpine
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o lambda-handler
# FROM public.ecr.aws/lambda/provided:al2023
FROM alpine:3.9
COPY --from=build-image /src/lambda-handler .
ENTRYPOINT ./lambda-handler
