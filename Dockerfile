# syntax=docker/dockerfile:experimental
# Build container
FROM --platform=$TARGETPLATFORM golang AS build
ARG DIBS_TARGET
ARG TARGETPLATFORM

WORKDIR /app

RUN curl -Lo /tmp/dibs https://nx904.your-storageshare.de/s/ZWxkmmQW37fHt9J/download
RUN install /tmp/dibs /usr/local/bin

ADD . .

RUN dibs -generateSources
RUN dibs -build

# Run container
FROM --platform=$TARGETPLATFORM alpine
ARG DIBS_TARGET
ARG TARGETPLATFORM

COPY --from=build /app/.bin/binaries/dibs* /usr/local/bin/dibs

CMD /usr/local/bin/dibs
