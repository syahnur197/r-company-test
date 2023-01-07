FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod vendor
COPY . ./
RUN apk --no-cache add curl
RUN go build -o /rakuten-app
EXPOSE 4000
CMD [ "/rakuten-app" ]