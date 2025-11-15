include .env
export $(shell sed 's/=.*//' .env)

# Variables
DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}

# Tasks (also called "targets")
migrate-up:
	@migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@migrate -path migrations -database "$(DB_URL)" down 1
	
migrate-drop:
	@migrate -path migrations -database "$(DB_URL)" drop

migrate-version:
	@migrate -path migrations -database "$(DB_URL)" version