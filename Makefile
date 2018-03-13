linux:
	GOOS=linux GOARCH=amd64 go build serv.go

mac:
	GOOS=darwin GOARCH=amd64 go build serv.go

clean:
	rm -f serv

cert:
	openssl ecparam -genkey -name secp384r1 -out serv.key && chmod 600 serv.key
	openssl req -new -x509 -sha256 -key serv.key -out serv.crt -days 3650 && chmod 600 serv.crt
