FROM golang:1.14 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/maestro ./cmd/maestro

FROM alpine
COPY --from=build /usr/local/bin/maestro /bin/maestro
ENTRYPOINT ["/bin/maestro"]