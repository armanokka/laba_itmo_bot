FROM golang:1.19.1-alpine3.16
WORKDIR app

ENV TRANSLOBOT_TOKEN "1934369237:AAGzGrSPC8hOf6suvJEv_fbC8lxqhqHrEs4"
ENV TRANSLOBOT_DASHBOT_TOKEN "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
ENV TRANSLOBOT_DSN "host=158.160.44.238 user=translobot password=i72aam8y#?gHduYA48ThJ741koEzQm4kOaqO5nt8AtFFxk9QjZCGB3{UK16Sc%Bd dbname=translobot port=5432 TimeZone=Europe/Moscow"
ENV TRANSLOBOT_DEBUG true
COPY ./ ./
RUN go mod tidy
RUN go build .
ENTRYPOINT [ "./translobot" ]
