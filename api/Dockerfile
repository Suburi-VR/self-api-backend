FROM golang:1.16.14

RUN  go get -u github.com/gin-gonic/gin

WORKDIR /usr/api

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

ENV AWS_ACCESS_KEY_ID=AKIA3OBNHHXY3UDDCG5Y
ENV AWS_SECRET_ACCESS_KEY=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
ENV REGION=ap-northeast-1

COPY . .

CMD ["go", "run", "main.go"]