FROM golang:1.14 AS build

ARG version=unknown
ARG build=unknown
ARG label=unknown

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /usr/local/bin/semaphore -ldflags "-X main.version=${version} -X main.build=${build} -X main.label=${label}" ./cmd/semaphore

FROM alpine
COPY --from=build /usr/local/bin/semaphore /bin/semaphore

RUN mkdir -p /etc/semaphore/
WORKDIR /etc/semaphore

ENTRYPOINT ["/bin/semaphore", "daemon"]