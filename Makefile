# Variabel
BINARY_NAME=main
MAIN_PATH=./cmd/api/main.go

# --- MODE LOKAL (Hybrid) ---
# Gunakan ini saat develop kode agar cepat (DB di Podman, App di Terminal)
dev-local:
	@echo "Menjalankan Karima Store di port 8080..."
	APP_PORT=8080 DB_HOST=localhost REDIS_HOST=localhost REDIS_PORT=6380 go run $(MAIN_PATH)

# --- MODE PODMAN (Full Container) ---
# Bangun image dan jalankan semua layanan di kontainer
docker-up:
	@echo "Membangun dan menjalankan semua layanan di Podman..."
	podman-compose up -d --build

# Matikan semua layanan
docker-down:
	@echo "Menghentikan semua layanan..."
	podman-compose down

# --- KRATOS AUTH ---
# Jalankan Kratos Services
kratos-up:
	@echo "Menjalankan layanan Ory Kratos..."
	podman-compose -f docker-compose.yml -f docker-compose.kratos.yml up -d kratos-migrate kratos kratos-selfservice-ui-node mailslurper

# Matikan Kratos Services
kratos-down:
	@echo "Menghentikan layanan Ory Kratos..."
	podman-compose -f docker-compose.yml -f docker-compose.kratos.yml stop kratos-migrate kratos kratos-selfservice-ui-node mailslurper

# Lihat log aplikasi backend saja
logs:
	podman logs -f karima_store_backend

# --- UTILITY ---
# Merapikan library Go
tidy:
	go mod tidy
	go mod verify

# Masuk ke terminal database postgres
db-shell:
	podman exec -it karima_postgres psql -U karima_store -d karima_db

# Bersihkan image sampah (<none>)
clean:
	podman image prune -f

# --- SWAGGER ---
# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	$(HOME)/go/bin/swag init -g cmd/api/main.go