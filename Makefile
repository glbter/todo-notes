cont_app:
	docker build -t todo-note .

cont_db:
	docker build -t todo-postgres ./docker/db/


run:cont_app
	docker-compose up -d

all-migrate:
	cd internal/repo/postgres/migrations
		#&& tern migrate  #--destination 1

app-with-data: run all-migrate
	docker exec -it todonote_db_1 psql postgresql://postgres:123@localhost:5432/todoNote -c "$(cat ./dump.sql)"
	#bash test-set-up.sh

db-test: app-with-data
	/usr/local/go/bin/go test -cover -tags postgres_test ./...
	docker exec -it todonote_db_1 psql postgresql://postgres:123@localhost:5432/todoNote -c "truncate table users cascade;"

	#bash test-clean-up.sh
	docker-compose stop

integration-test: app-with-data
	go test -cover -tags postgres_test ./...
	docker exec -it todonote_db_1 psql postgresql://postgres:123@localhost:5432/todoNote -c "truncate table users cascade;"
	#bash test-clean-up.sh
	docker-compose stop

#make -f Makefile test


start-postgres:
	docker run -d -p 5432:5432 -e POSTGRES_DB=todoNote -e POSTGRES_PASSWORD=123 -e PGDATA=/var/lib/postgresql/data -v /postgres/todoNote:/var/lib/postgresql/data postgres

# apm server
apm-start:
	docker run -d -p 8200:8200 --name=apm-server --user=apm-server docker.elastic.co/apm/apm-server:7.15.0 --strict.perms=false -e -E output.elasticsearch.hosts=["host.docker.internal:9200"]
	docker run -d --name es01-test --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.15.1
	docker run -d --name kib01-test --net elastic -p 5601:5601 -e "ELASTICSEARCH_HOSTS=http://es01-test:9200" docker.elastic.co/kibana/kibana:7.15.1

apm-stop:
	docker container stop es01-test
	docker container stop kib01-test
	docker container stop apm-server

	docker container rm kib01-test
	docker container rm es01-test
	docker container rm apm-server


