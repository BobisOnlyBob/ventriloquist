FROM xena/go:1.10 AS build
COPY . /root/go/src/github.com/withinsoft/ventriloquist
RUN GOBIN=/usr/local/bin go install github.com/withinsoft/ventriloquist/cmd/ventriloquist

FROM xena/alpine
COPY --from=build /usr/local/bin/ventriloquist /usr/local/bin/ventriloquist
VOLUME /data
ENV DB_PATH /data/tulpas.db
CMD ["/usr/local/bin/ventriloquist"]
