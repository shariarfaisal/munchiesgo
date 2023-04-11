postgres:
	docker run --name postgres14 -p 6500:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:14.6-alpine
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root munchies_backend

dropdb:
	docker exec -it postgres14 dropdb munchies_backend

newmigration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "postgresql://root:admin@localhost:6500/munchies_backend?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:admin@localhost:6500/munchies_backend?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

	