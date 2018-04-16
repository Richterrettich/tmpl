FROM golang:1.10.1-alpine3.7
WORKDIR /go/src/tmpl
COPY . .
RUN go build -o tmpl -ldflags="-s -w"


FROM scratch
COPY --from=0 /go/src/tmpl/tmpl /tmpl
ENTRYPOINT [ "/tmpl" ]