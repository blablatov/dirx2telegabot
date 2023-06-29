# Multi-stage dirx2telegabot build
# Многоэтапная сборка dirx2telegabot

FROM golang AS build

ENV location /go/src/github.com/blablatov/dirx2telegabot

WORKDIR ${location}/dirx2telegabot

ADD ./dirx2telegabot.go ${location}/dirx2telegabot

RUN go mod init github.com/blablatov/dirx2telegabot/dirx2telegabot

RUN CGO_ENABLED=0 go build -o dirx2telegabot

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.

FROM scratch
COPY --from=build ./dirx2telegabot ./dirx2telegabot

ENTRYPOINT ["./dirx2telegabot"]
EXPOSE 50051