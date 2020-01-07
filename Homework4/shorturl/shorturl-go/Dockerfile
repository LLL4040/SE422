FROM golang:latest

ADD  ./github.com /go/src/github.com
ADD ./ /home/app/shorturl

RUN cd /home/app/shorturl && go build shorturl.go

EXPOSE 9000

ENTRYPOINT [ "/home/app/shorturl/shorturl" ]