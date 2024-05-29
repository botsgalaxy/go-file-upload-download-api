FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=builder /go/src/app/app .

COPY .env .

RUN chmod +x ./app

CMD ["./app"]
