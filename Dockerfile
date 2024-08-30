FROM golang:1.23.0-alpine

WORKDIR ./app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /kamibooking

EXPOSE 3000

CMD ["/kamibooking"]