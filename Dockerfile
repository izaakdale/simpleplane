FROM golang:1.20-alpine as builder
WORKDIR /

COPY . .
RUN go mod download

RUN go build -o simpleplane .

FROM scratch
WORKDIR /bin
COPY --from=builder /simpleplane /bin

CMD [ "/bin/simpleplane" ]