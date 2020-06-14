# Sabbatical
 some stuff i want to explore, post or mid sabbatical
 
 **please note:** this is not to aimed to be a feature complete / production ready product at any point,
 but rather a personal playground for microservices. thereby things like user input validation,
 is not going to get implemented in the API

so far:
 - Basic versioned Rest/CRUD API written in golang
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
