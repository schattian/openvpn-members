FROM golang:alpine AS builder

# RUN apk add git

LABEL maintainer="Sebastián Chamena <sebachamena@gmail.com>"

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app .

FROM scratch

WORKDIR /app

COPY --from=builder /app/app ./app

CMD ["/app/app"]