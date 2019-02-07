DOMAIN_NAME=localhost
CA_NAME=ca
SERVER_NAME=server
CLIENT_NAME=client

all: cert protobuf

cert: server-cert ca-cert client-cert

ca-cert: server-cert
	openssl genrsa -out $(CA_NAME).key 2048
	openssl req -new -key $(CA_NAME).key -out $(CA_NAME).csr -subj "/CN=$(CA_NAME)"

	openssl req -x509 -sha256 -nodes -days 365 -key $(CA_NAME).key \
	-out $(CA_NAME).crt -subj "/CN=Certified" # self sign certificate

	openssl x509 -req -days 365 -in ./server/$(SERVER_NAME).csr \
	-CA $(CA_NAME).crt -CAkey $(CA_NAME).key -CAcreateserial -out $(SERVER_NAME).crt
	mv $(CA_NAME).* ./ca
	mv $(SERVER_NAME).crt ./server

server-cert:
	openssl req -new \
	-newkey rsa:2048 -nodes -keyout $(SERVER_NAME).key \
	-out $(SERVER_NAME).csr \
	-subj "/CN=$(DOMAIN_NAME)"
	mv $(SERVER_NAME).* ./server

client-cert:
	cp ./ca/$(CA_NAME).crt ./client
	openssl req -new \
	-newkey rsa:2048 -nodes -keyout $(CLIENT_NAME).key \
	-out $(CLIENT_NAME).csr \
	-subj "/CN=$(DOMAIN_NAME)"
	mv $(CLIENT_NAME).* ./client

protobuf: dummy.proto
	protoc -I. dummy.proto --go_out=plugins=grpc:./server --go_out=plugins=grpc:./client


.PHONY : clean
clean:
	rm -f .*/*/*.pb.go
	rm -f .*/*/$(CA_NAME).*
	rm -f .*/*/$(SERVER_NAME).*
	rm -f .*/*/$(CLIENT_NAME).*
