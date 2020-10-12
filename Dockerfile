FROM golang:1.15.2 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -json
#RUN go mod graph | cut -d ' ' -f 2 | sort | uniq | tr '\n' ' ' | xargs go get -v

COPY . .
RUN find . -type f && CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build app/app .
CMD ["./app"]
