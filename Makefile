.PHONY: run-db
run-db:
	$(info #Run postgresql service and start migration)
	docker-compose up -d && psql -U user -h localhost -f migration.sql
