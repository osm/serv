Serve a single static file on a port that suits you.

Example
```sh
# Generate a self signed key and certificate
make cert

# Create a file with the basic auth credentials
echo -n "foo:bar" | base64 > credentials.txt && chmod 600 credentials.txt

# Build the server
make linux

# Start the server
./serv -cred credentials.txt -key serv.key -cert serv.crt -file serv.go

# Test it
curl -u foo:bar -k https://localhost:4554
```
