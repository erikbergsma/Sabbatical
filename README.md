# Sabbatical
 some stuff i want to explore, post or mid sabbatical

so far:
 - CRUD API written in golang + redis as a DB
 - API can run against an external Redis service, like AWS' Elasticache
 - API frontend in angularjs (under NGINX)
 - frontend, API and DB are seperate containerized services
 - docker-compose for local building/running/testing
 - github action for building && deploying to AWS ECS
 
## to see this project in action
```shell
cd Sabbatical/src/web/frontend
sudo docker-compose up
```
docker-compose will then either run or build all the containers needed, pretty cool
