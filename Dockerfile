FROM golang:1.22.3-alpine3.20

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

# Install dbmate
RUN apk add --no-cache curl bash && \
    curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 -o /usr/local/bin/dbmate && \
    chmod +x /usr/local/bin/dbmate

COPY . .

RUN go mod tidy

CMD ["sh", "-c", "dbmate up"]

CMD ["air", "-c", ".air.toml"]
