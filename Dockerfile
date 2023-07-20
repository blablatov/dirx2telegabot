# Yandex Compute Cloud:
# $ sudo docker build . -t cr.yandex/${REGISTRY_ID}/debian:dirx2telegabot -f Dockerfile
# $ sudo docker run --name dirx2telegabot -p 8077:8077 -d cr.yandex/${REGISTRY_ID}/debian:dirx2telegabot
# test:
# $ curl external_address_vm:8077/Уведомление

# Local:
# $ docker build -t dirx2telegabot -f Dockerfile
# $ docker run --name dirx2telegabot -p 8077:8077 -d dirx2telegabot
# test:
# $ curl localhost:8077/Уведомление

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dirx2telegabot

EXPOSE 8077

CMD ["/dirx2telegabot"]