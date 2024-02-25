FROM alpine

WORKDIR goauth
COPY goauth .

CMD ["./goauth"]
EXPOSE 8080
