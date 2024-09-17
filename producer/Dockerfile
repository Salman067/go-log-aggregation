FROM golang:1.18-alpine
ENV GO111MODULE=on

RUN mkdir /app
WORKDIR /app

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/producer

RUN go build -o /main .

EXPOSE 4000

CMD ["/main"]
