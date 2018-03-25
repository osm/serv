FROM golang:latest
WORKDIR /app
COPY . /app
RUN make linux
CMD ["./serv", "-key", "serv.key", "-cert", "serv.crt", "-cred", "credentials.txt"]
