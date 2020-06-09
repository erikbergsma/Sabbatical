# Sabbatical
 some stuff i want to explore, post or mid sabbatical

so far:
 - Basic Rest/CRUD API written in golang
 - API can run against a local redis (container) or against an external Redis service, like AWS' Elasticache
 - API frontend in angularjs (under NGINX)
 - Frontend, API and DB are seperate containerized services
 - docker-compose for local building/running/testing
 - github action for building each service to amd64 and arm/v7 containers
 - arm/v7 support so that everything can run on a Raspberry pi!
 - tools to use letsencrypt and use/build a webhook server

## to see this project in action
```shell
cd Sabbatical/src/web/frontend
sudo docker-compose up
```
docker-compose will then either run or build all the containers needed, pretty cool
