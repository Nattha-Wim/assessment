FROM golang:1.19-alpine as build-base

WORKDIR /app1

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=unit -v ./...

RUN go build -o ./out/go-assessment .

# ====================


FROM alpine:3.16.2
COPY --from=build-base /app1/out/go-assessment /app1/go-assessment

CMD ["/app/go-assessment"]
