#!/bin/bash
script=$(cat ./dump.sql)
docker exec -it todonote_db_1 psql postgresql://postgres:123@localhost:5432/todoNote -c "$script"


	#docker cp .\dump.sql todonote_db_1:docker-entrypoint-initdb.d


#cat := $(if $(filter $(OS),Windows_NT),type,cat)
	#script := $(shell $(cat) ./dump.sql)
	# script = $(file < ./dump.sql)
#shell cat ./dump.sql #script = Get-Content .\dump.sql -Raw
