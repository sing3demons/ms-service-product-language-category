start: start-local mongo-init
start-docker: start-docker-up

start-local:
	docker compose -f docker-compose.local.yml up -d
	sleep 5 && echo "Waiting for mongo to start"
mongo-init:
	docker exec -it mongodb1 mongosh --eval "rs.initiate({_id:'my-replica-set', members:[{_id:0, host:'mongodb1:27017'}, {_id:1, host:'mongodb2:27018'}, {_id:2, host:'mongodb3:27019'}]})"
	docker exec -it mongodb1 mongosh --eval "rs.status()"
stop:
	docker-compose -f docker-compose.local.yml down
	sleep 5 && echo "Waiting for mongo to stop"
	rm -rf ./data
start-docker-up:
	docker compose -f docker-compose.db.yml up -d
	sleep 5 && echo "Waiting for mongo to start"
	docker exec -it mongodb1 mongosh --eval "rs.initiate({_id:'my-replica-set', members:[{_id:0, host:'mongodb1:27017'}, {_id:1, host:'mongodb2:27018'}, {_id:2, host:'mongodb3:27019'}]})"
	docker exec -it mongodb1 mongosh --eval "rs.status()"
	docker compose -f docker-compose.yml up -d

stop-docker:
	docker-compose -f docker-compose.yml down
	sleep 5 && echo "Waiting for mongo to stop"
	docker-compose -f docker-compose.db.yml down
	sleep 5 && echo "Waiting for mongo to stop"
	rm -rf ./data