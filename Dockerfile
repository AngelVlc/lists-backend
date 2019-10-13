FROM golang:1.12-alpine as base

RUN apk --update add git
RUN go get -u gopkg.in/mgo.v2 golang.org/x/crypto/bcrypt github.com/dgrijalva/jwt-go

ENV APP /go/src/github.com/AngelVlc/lists-backend
WORKDIR $APP

FROM base as test
RUN apk add gcc musl-dev
RUN go get -u github.com/stretchr/testify
COPY . $APP

FROM base as release
COPY . $APP
EXPOSE 5000
CMD ["go","run","."]
