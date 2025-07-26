db_version=postgres:latest
postgres_name=postgres
migrate_dir_path=../../internal/repositories/postgresql/migrations/
postgres_user=root
postgres_password=1234
db_name=database
db_port=5432
db_source=postgresql://$(postgres_user):$(postgres_password)@0.0.0.0:$(db_port)/$(db_name)?sslmode=disable
mg_create:
	@migrate create -ext sql -dir $(migrate_dir_path) -seq init_schema_$(NAME)

mgu:
	@migrate -path $(migrate_dir_path) -database $(db_source) -verbose up
mgd:
	@migrate -path $(migrate_dir_path) -database $(db_source) -verbose down
compose:
	cd ../../docker/ && docker-compose up -d --build

sqlc:
	cd ../sqlc && sqlc generate
