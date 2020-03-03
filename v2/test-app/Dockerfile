# Build container
FROM golang:alpine AS build

WORKDIR /app

RUN apk add -u curl

RUN curl -Lo /tmp/dibs https://nx904.your-storageshare.de/s/FkBmkPRPPC7Ea2A/download
RUN install /tmp/dibs /usr/local/bin

ADD . .

RUN dibs -build

# Run container
FROM alpine

COPY --from=build /app/.bin/binaries/test-app* /usr/local/bin

CMD /usr/local/bin/test-app
