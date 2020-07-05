FROM golang:1.14 AS build

ARG version=unknown
ARG build=unkown
ARG label=unkown

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /usr/local/bin/maestro -ldflags "-X main.version=${version} -X main.build=${build} -X main.label=${label}" ./cmd/maestro

FROM alpine
COPY --from=build /usr/local/bin/maestro /bin/maestro

RUN mkdir -p /etc/maestro/
COPY ./resources/default/ /etc/maestro/
WORKDIR /etc/maestro

ENTRYPOINT ["/bin/maestro", "daemon"]