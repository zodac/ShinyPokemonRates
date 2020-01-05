FROM golang
LABEL maintainer="zodac <arouge110@msn.com>"

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

EXPOSE 5000

CMD ["./main"]