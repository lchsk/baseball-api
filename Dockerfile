FROM golang

# RUN git clone https://github.com/lchsk/baseballapi /go/src/github.com/lchsk/baseballapi

ADD . /go/src/github.com/lchsk/baseballapi

RUN go get ./...

RUN go install github.com/lchsk/baseballapi

ENTRYPOINT ["/go/bin/baseballapi"]

EXPOSE 8000 5433