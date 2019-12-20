# syntax=docker/dockerfile:experimental
FROM --platform=$TARGETPLATFORM golang:1.13.5-buster AS build
WORKDIR /app
ARG TARGETPLATFORM

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./.dibs.yml ./.dibs.yml
COPY ./main.go ./main.go
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./.git ./.git

RUN go run main.go pipeline build assets

FROM --platform=$TARGETPLATFORM debian:buster-slim
ARG TARGETPLATFORM

COPY ./.dibs.yml ./.dibs.yml

COPY --from=build /app/.bin/dibs-* /usr/local/bin/dibs

EXPOSE 32000

CMD dibs pipeline sync server
