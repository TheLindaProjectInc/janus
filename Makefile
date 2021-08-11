ifndef GOBIN
GOBIN := $(GOPATH)/bin
endif

ifdef JANUS_PORT
JANUS_PORT := $(JANUS_PORT)
else
JANUS_PORT := 23889
endif

check-env:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

.PHONY: install
install: 
	go install github.com/TheLindaProjectInc/janus/cli/janus

.PHONY: release
release: darwin linux

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/TheLindaProjectInc/janus/cli/janus

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/janus-linux-amd64 github.com/TheLindaProjectInc/janus/cli/janus

.PHONY: quick-start
quick-start:
	cd docker && ./spin_up.sh && cd ..

.PHONY: docker-dev
docker-dev:
	docker build -t metrixcoin/janus:dev .
	
.PHONY: local-dev
local-dev: check-env
	go install github.com/TheLindaProjectInc/janus/cli/janus
	docker run --rm --name metrix_testchain -d -p 33851:33851 metrixcoin/metrix metrixd -regtest -rpcbind=0.0.0.0:33851 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=metrix -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/TheLindaProjectInc/janus/docker/fill_user_account.sh metrix_testchain:.
	docker exec metrix_testchain /bin/sh -c ./fill_user_account.sh
	METRIX_RPC=http://metrix:testpasswd@localhost:33851 METRIX_NETWORK=regtest $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev

.PHONY: local-dev-https
local-dev-https: check-env
	go install github.com/TheLindaProjectInc/janus/cli/janus
	docker run --rm --name metrix_testchain -d -p 33851:33851 metrixcoin/metrix metrixd -regtest -rpcbind=0.0.0.0:33851 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=metrix -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/TheLindaProjectInc/janus/docker/fill_user_account.sh metrix_testchain:.
	docker exec metrix_testchain /bin/sh -c ./fill_user_account.sh > /dev/null&
	METRIX_RPC=http://metrix:testpasswd@localhost:33851 METRIX_NETWORK=regtest $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev --https-key https/key.pem --https-cert https/cert.pem

.PHONY: local-dev-logs
local-dev-logs: check-env
	go install github.com/TheLindaProjectInc/janus/cli/janus
	docker run --rm --name metrix_testchain -d -p 33851:33851 metrixcoin/metrix:dev metrixd -regtest -rpcbind=0.0.0.0:33851 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=metrix -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/TheLindaProjectInc/janus/docker/fill_user_account.sh metrix_testchain:.
	docker exec metrix_testchain /bin/sh -c ./fill_user_account.sh
	METRIX_RPC=http://metrix:testpasswd@localhost:33851 METRIX_NETWORK=regtest $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev > janus_dev_logs.txt

.PHONY: unit-tests
unit-tests: check-env
	go test -v ./... -timeout 30s

docker-build-unit-tests:
	docker build -t metrixcoin/tests.janus -f ./docker/unittests.Dockerfile .

docker-unit-tests:
	docker run --rm -v `pwd`:/go/src/github.com/TheLindaProjectInc/janus metrixcoin/tests.janus

docker-tests: docker-build-unit-tests docker-unit-tests openzeppelin-docker-compose

docker-configure-https: docker-configure-https-build
	docker/setup_self_signed_https.sh

docker-configure-https-build:
	docker build -t metrixcoin/openssl.janus -f ./docker/openssl.Dockerfile ./docker

# -------------------------------------------------------------------------------------------------------------------
# NOTE:
# 	The following make rules are only for local test purposes
# 
# 	Both run-janus and run-metrix must be invoked. Invocation order may be independent, 
# 	however it's much simpler to do in the following order: 
# 		(1) make run-metrix 
# 			To stop Metrix node you should invoke: make stop-metrix
# 		(2) make run-janus
# 			To stop Janus service just press Ctrl + C in the running terminal

# Runs current Janus implementation
run-janus:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/cli/janus/main.go \
		--metrix-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:33851 \
		--metrix-network=regtest \
		--bind=0.0.0.0 \
		--port=23888 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--dev

run-janus-https:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/cli/janus/main.go \
		--metrix-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:33851 \
		--metrix-network=regtest \
		--bind=0.0.0.0 \
		--port=23888 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--dev \
		--https-key https/key.pem \
		--https-cert https/cert.pem

test_user = metrix
test_user_passwd = testpasswd

# Runs docker container of metrix locally and starts metrixd inside of it
run-metrix:
	@ printf "\nRunning metrix...\n\n"
	@ printf "\n(1) Starting container...\n\n"
	docker run ${metrix_container_flags} metrixcoin/metrix metrixd ${metrixd_flags} > /dev/null

	@ printf "\n(2) Importing test accounts...\n\n"
	@ sleep 3
	docker cp ${shell pwd}/docker/fill_user_account.sh ${metrix_container_name}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${metrix_container_name} /bin/sh -c ./fill_user_account.sh > /dev/null
	@ printf "\n... Done\n\n"

metrix_container_name = test-chain

# TODO: Research -v
metrix_container_flags = \
	--rm -d \
	--name ${metrix_container_name} \
	-v ${shell pwd}/dapp \
	-p 33851:33851

# TODO: research flags
metrixd_flags = \
	-regtest \
	-rpcbind=0.0.0.0:33851 \
	-rpcallowip=0.0.0.0/0 \
	-logevents \
	-addrindex \
	-reindex \
	-txindex \
	-rpcuser=${test_user} \
	-rpcpassword=${test_user_passwd} \
	-deprecatedrpc=accounts \
	-printtoconsole

# Starts continuously printing Metrix container logs to the invoking terminal
follow-metrix-logs:
	@ printf "\nFollowing metrix logs...\n\n"
		docker logs -f ${metrix_container_name}

open-metrix-bash:
	@ printf "\nOpening metrix bash...\n\n"
		docker exec -it ${metrix_container_name} bash

# Stops docker container of metrix
stop-metrix:
	@ printf "\nStopping metrix...\n\n"
		docker kill `docker container ps | grep ${metrix_container_name} | cut -d ' ' -f1` > /dev/null
	@ printf "\n... Done\n\n"

restart-metrix: stop-metrix run-metrix

submodules:
	git submodules init

# Run openzeppelin tests, Janus/METRIX needs to already be running
openzeppelin:
	cd testing && make openzeppelin

# Run openzeppelin tests in docker
# Janus and METRIX need to already be running
openzeppelin-docker:
	cd testing && make openzeppelin-docker

# Run openzeppelin tests in docker-compose
openzeppelin-docker-compose:
	cd testing && make openzeppelin-docker-compose