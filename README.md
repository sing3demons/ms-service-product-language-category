# ms-service-product-language-category
golang mongo kafka nodejs


```start-db
docker-compose up -d
docker exec -it mongodb1 mongosh --eval "rs.initiate({_id:\"my-replica-set\",members:[{_id:0,host:\"mongodb1:27017\"},{_id:1,host:\"mongodb2:27018\"},{_id:2,host:\"mongodb3:27019\"}]})"
```

```windows
docker exec -it mongodb1 mongosh --eval "rs.initiate({_id:'my-replica-set', members:[{_id:0, host:'mongodb1:27017'}, {_id:1, host:'mongodb2:27018'}, {_id:2, host:'mongodb3:27019'}]})"

```