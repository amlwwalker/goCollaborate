FROM base-git2go

ADD ./compiler /go/src/github.com/compiler

WORKDIR /go/src/github.com/compiler

#install...
RUN go get github.com/gorilla/websocket

RUN go get github.com/fatih/color

#get fsnotify
RUN go get github.com/howeyc/fsnotify

RUN go build .

RUN ls -l ./

CMD /go/src/github.com/compiler/compiler -directory=/linked -command=/linked/linked
EXPOSE 8080
