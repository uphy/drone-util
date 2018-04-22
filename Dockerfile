FROM golang:1.9 as builder

WORKDIR /go/src/github.com/uphy/drone-util
COPY . /go/src/github.com/uphy/drone-util
RUN CGO_ENABLED=0 go build -o /drone-util

FROM alpine:3.7

COPY --from=builder /drone-util /bin/drone-util
RUN chmod +x /bin/drone-util
ENTRYPOINT [ "/bin/drone-util" ]
CMD [ "export" ]
