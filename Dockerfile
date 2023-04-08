FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN mkdir logs
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

CMD ["/main"]