FROM golang:1.16.14

RUN  go get -u github.com/gin-gonic/gin

WORKDIR /usr/selfapi

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

ENV AWS_ACCESS_KEY_ID=hironori1
ENV AWS_SECRET_ACCESS_KEY=hironori1
ENV REGION=ap-northeast-1

EXPOSE 8080

COPY . .

CMD ["go", "run", "user.go"]