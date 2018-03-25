# Serve a single static file on a port that suits you.

## Example

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

## Run from Docker

```sh
# Generate a self signed key and certificate
make cert

# Create a file with the basic auth credentials
echo -n "foo:bar" | base64 > credentials.txt && chmod 600 credentials.txt

# Build the Docker image
docker build -t osm-serv .

# Create a file to host
mkdir -p /tmp/serv_storage
echo foo > /tmp/serv_storage/hosted_file.txt

# Run the docker image
docker run \
	-p "4554:4554/tcp" \
	-e "FILE=/tmp/hosted_file.txt" \
	-v "/tmp/serv_storage:/tmp" \
	osm-serv

# Test it
curl -u foo:bar -k https://localhost:4554
