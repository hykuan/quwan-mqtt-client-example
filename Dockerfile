# build stage
FROM golang:alpine AS build-env

RUN apk update && \
    apk --update add --no-cache git openssh ca-certificates curl bash make build-base && \
    rm -rf /var/cache/apk/* && rm -rf /var/lib/apt/lists/*

ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# final stage
FROM alpine:3.8

WORKDIR /app
COPY --from=build-env /src/app /app/

EXPOSE 8080

CMD ["./app"]