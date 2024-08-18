CONFIG_PATH=${HOME}/Github/distributed_services_with_go/.certs

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl/bin/cfssl gencert \
		-initca test/ca-csr.json | cfssl/bin/cfssljson -bare ca

	cfssl/bin/cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssl/bin/cfssljson -bare server

	cfssl/bin/cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		test/client-csr.json | cfssl/bin/cfssljson -bare client
	
	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: test
test:
	go test -race ./...

.PHONY: compile
compile:
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=. \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative
