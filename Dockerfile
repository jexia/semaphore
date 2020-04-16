FROM golang:1.14 AS build

WORKDIR /app
COPY cmd/maestro/go.mod cmd/maestro/go.sum ./
RUN go mod download

COPY cmd/maestro .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/maestro .

FROM scratch
COPY --from=build /usr/local/bin/maestro /maestro
ENTRYPOINT ["/maestro"]