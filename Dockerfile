FROM golang:1.22-alpine AS build-env

WORKDIR /tmp/build
ADD . /tmp/build
# -ldlflags '-s' to strip binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -ldflags '-w -s'

###

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/pcvolkmer/rwdp-obds-cleanup"
LABEL org.opencontainers.image.licenses="AGPLv3"
LABEL org.opencontainers.image.description="Simple pipeline step to clean up incoming kafka messages"

COPY --from=build-env /tmp/build/app /rwdp-obds-cleanup

USER 8000:8000

ENTRYPOINT ["/rwdp-obds-cleanup"]