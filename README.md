# Sabbatical
 some stuff i want to explore, post or mid sabbatical

so far:
 - CRUD API written in golang, Redis as a DB
 - API frontend in angularjs (under NGINX)
 - API can run against an external Redis service, like AWS' Elasticache
 - frontend, API and DB are seperate containerized services
 - docker-compose for local building/running/testing
 - github action for building each service to amd64 and arm/v7 containers
 - arm/v7 support so that everything can run on a Raspberry pi!

## to see this project in action
```shell
cd Sabbatical/src/web/frontend
sudo docker-compose up
```
docker-compose will then either run or build all the containers needed, pretty cool
