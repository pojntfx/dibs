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

RUN go run main.go pipeline build assets --platform $TARGETPLATFORM

FROM --platform=$TARGETPLATFORM debian:buster-slim
ARG TARGETPLATFORM

COPY ./.dibs.yml ./.dibs.yml

COPY --from=build /app/.bin/godibs-* /usr/local/bin/godibs

EXPOSE 25000
CMD /usr/local/bin/godibs server
