DB_URL=postgresql://root:root@localhost:5432/movies?sslmode=disable

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:12-alpine

server:
	go run cmd/main.go

proto:
	rm -rf domain/pb/*.go
	protoc --proto_path=domain/proto --go_out=domain/pb --go_opt=paths=source_relative \
    --go-grpc_out=domain/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=domain/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=movies \
    domain/proto/*.proto

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root movies

dropdb:
	docker exec -it postgres12 dropdb movies

sqlc:
	sqlc generate

migrateup:
	migrate -path migrations -database "$(DB_URL)" --verbose up

createmigrate:
	migrate create -ext sql -dir migrations -seq init_schema

migratedown:
	migrate -path migrations -database "$(DB_URL)" --verbose down

newmigration:
	migrate create -ext sql -dir migrations -seq $(name)
	
test:
	go test -v -cover ./...

mock:
	mockgen -package mock -destination internal/mock/store.go github.com/lovelyoyrmia/movies-api/internal/db Store

seed:
	go run seeds/seed.go