FROM golang:1.14 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/maestro .

FROM scratch
COPY --from=build /usr/local/bin/maestro /maestro
ENTRYPOINT ["/maestro"]