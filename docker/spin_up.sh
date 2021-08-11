#!/bin/sh
docker-compose -f  ${GOPATH}/src/github.com/TheLindaProjectInc/janus/docker/quick_start/docker-compose.yml up -d 
sleep 3 #executing too fast causes some errors 
docker cp ${GOPATH}/src/github.com/TheLindaProjectInc/janus/docker/fill_user_account.sh metrix_testchain:.
docker exec metrix_testchain /bin/sh -c ./fill_user_account.sh