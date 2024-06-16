build:
	docker-compose build library_app

run:
	docker-compose up library_app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
	#migrate -path ./schema -database 'postgres://postgres:qwerty@db:5432/postgres?sslmode=disable' up