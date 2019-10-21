FROM golang:1.12 as base

ENV APP /go/src/github.com/AngelVlc/lists-backend
WORKDIR $APP
COPY . $APP

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

FROM base as test
RUN go get -u github.com/stretchr/testify

FROM alpine as release
COPY --from=base /go/bin/app /
CMD [ "./app" ]
