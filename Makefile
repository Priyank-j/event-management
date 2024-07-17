.PHONY: build run logs dockerstop run-worker build-worker dockerstopworker
.SILENT: build run logs dockerstop run-worker build-worker dockerstopworker

GH_PAT ?= $(shell bash -c 'read -s -p "Github Personal Access Token: " pwd; echo $$pwd')

ifeq ($(stage), prod)
 cloud-watch-logs= --log-driver=awslogs --log-opt awslogs-group=lending-middleware-prod-logs
else
 cloud-watch-logs=
endif


init:
	docker network create finboxnet

build:
	docker build -t go-docker-event-optimized-${stage} --build-arg="GH_PAT=${GH_PAT}" .	

logs:
	docker logs $(shell docker ps | grep 'go-docker-event-optimized-apis-${stage}' | awk '{ print $$1 }') -f ;


dockerstop:
	did=$(shell docker ps | grep 'go-docker-event-optimized-apis-${stage}' | awk '{ print $$1 }'); \
	if [ "$$did" ]; \
	then docker stop $$did; \
	docker logs $$did > ../logs/"$$did"_"$(shell date +"%Y_%m_%d_%I_%M_%p").log" ; \
	echo "go-docker-event-optimized-${stage} stopped"; \
	docker rm 'go-docker-event-optimized-apis-${stage}'; \
	echo "go-docker-event-optimized-${stage} removed"; \
 	else echo "no go-docker-event-optimized-${stage} container found"; fi; 

run: dockerstop
	port=$(shell echo ${stage} | sed 's/[^0-9]*//g' | sed 's/^$$/1/' | awk '{print $$1 "+3332"}' | bc); \
	if [ ${stage} = "local" ]; then docker run -d --name go-docker-event-optimized-apis-${stage} -e STAGE=${stage} -e AWS_PROFILE=${aws_profile} -v ${HOME}/.aws/credentials:/root/.aws/credentials:ro -p $$port:3332 -v ~/logs/go-docker:/app/logs go-docker-event-optimized-${stage};\
	else docker run -d --restart unless-stopped --network finboxnet --name go-docker-event-optimized-apis-${stage} --pids-limit 50 -e STAGE=${stage} -e DD_ENV=${stage} -e DD_SERVICE="lending-middleware-apis" -e DD_VERSION=$(shell git rev-parse --short HEAD) -l com.datadoghq.tags.env=${stage} -l com.datadoghq.tags.service="lending-middleware-apis" -l com.datadoghq.tags.version=$(shell git rev-parse --short HEAD) -p $$port:3332 -v ~/logs/go-docker:/app/logs ${cloud-watch-logs} go-docker-event-optimized-${stage}; fi;\


all: build run logs

