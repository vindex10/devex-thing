FROM golang:latest AS tmp

COPY ./ /root/dist/
WORKDIR /root/dist
RUN go build internal/main.go

FROM golang:latest

COPY --from=tmp /root/dist/main /root/main
WORKDIR /root
EXPOSE 8080
CMD ["./main"]
