#!/bin/bash
docker exec -it todonote_db_1 psql postgresql://postgres:123@localhost:5432/todoNote -c "truncate table users cascade;"
