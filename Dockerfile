FROM golang:1.20

RUN git clone https://github.com/blablatov/dirx2telegabot.git
WORKDIR dirx2telegabot

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dirx2telegabot

EXPOSE 8077

CMD ["/dirx2telegabot"]