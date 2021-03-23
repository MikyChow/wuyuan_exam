# wuyuan_exam
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
### Setup this project
```
# deploy docker-compose
docker-compose up -d 

# init database
CREATE TABLE "tasks" (
  "id" bigint NOT NULL Primary Key,
  "parent_id" bigint NOT NULL,
  "id_path" character varying(300) NOT NULL,
  "status" smallint NOT NULL,
  "description" text NULL,
  "duedate" date NOT NULL
);