# Migration Commands
c_m:
	@echo "Creating migrations..."
	migrate create -ext sql -dir db/migrations $(name)

# PostgreSQL Commands
p_up:
	@echo "Starting PostgreSQL..."
	docker compose up -d
p_down:
	@echo "Stopping and removing PostgreSQL..."
	docker compose down
db_up:
	@echo "Creating a database..."
	docker exec -it ecommerce_postgres createdb --username=root --owner=root commerce_db
	docker exec -it ecommerce_postgres_live createdb --username=root --owner=root livedb

db_down:
	@echo "Dropping the database..."
	docker exec -it ecommerce_postgres dropdb --username=root commerce_db
	docker exec -it ecommerce_postgres_live dropdb --username=root livedb

# Migration Up and Down Commands
m_up:
	@echo "Applying database migrations (up)..."
	migrate -path db/migrations -database "postgres://root:secret@localhost:5433/commerce_db?sslmode=disable" up
	migrate -path db/migrations -database "postgres://root:secret@localhost:5434/livedb?sslmode=disable" up

m_down:
	@echo "Reverting database migrations (down)..."
	migrate -path db/migrations -database "postgres://root:secret@localhost:5433/commerce_db?sslmode=disable" down
	migrate -path db/migrations -database "postgres://root:secret@localhost:5434/livedb?sslmode=disable" down
# SQLC Commands
sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

# Testing
test:
	@echo "Running tests..."
	go test -v -cover ./...

# package for JWT token
jwt:
	go get github.com/golang-jwt/jwt
# Start Development Server
start:
	@echo "Starting the development server..."
	CompileDaemon -command="./ecommerce_backend"

