FROM golang:1.19 AS builder
WORKDIR /go/src/celso
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /go/src/celso/app ./
CMD ["./app"]