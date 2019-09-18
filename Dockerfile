FROM golang:1.12-alpine

RUN apk --update add git

RUN go get -u gopkg.in/mgo.v2

ENV APP /go/src/github.com/AngelVlc/lists-backend
WORKDIR $APP

COPY . $APP

EXPOSE 5000

CMD ["go","run","."]
