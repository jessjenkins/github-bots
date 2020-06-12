FROM golang:1.14.4 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -json
COPY . .
RUN find . -type f && CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build app/app .
CMD ["./app"]