FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/mediocregopher/radix.v2/redis"
RUN go build -o media 
CMD ["/app/main"]
EXPOSE 3006