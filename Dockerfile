FROM alpine:latest
WORKDIR app
COPY ./laba_itmo_bot ./laba_itmo_bot
ENTRYPOINT [ "./laba_itmo_bot" ]
