FROM golang:1.12.5-stretch

COPY bin/counter /counter
COPY . .

EXPOSE 8888 8888
ENTRYPOINT ["/counter"]