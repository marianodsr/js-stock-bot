FROM golang:alpine
WORKDIR /go/src/stock-bot
COPY . .
RUN go build -o /go/bin/stock-bot .

CMD [ "/go/bin/stock-bot" ]