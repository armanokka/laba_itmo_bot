FROM golang:1.19-alpine3.16
WORKDIR app
COPY ./laba_itmo_bot ./
ENTRYPOINT [ "./laba_itmo_bot" ]
