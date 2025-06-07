
# Makefile

# Run all services
up:
	docker compose up --build

# Stop all containers
down:
	docker compose down

# Rebuild only frontend
build-frontend:
	docker compose build frontend

# Rebuild only backend
build-backend:
	docker compose build backend

# Tail logs
logs:
	docker compose logs -f

# Restart everything
restart:
	docker compose down && docker compose up --build

# Clean all unused Docker resources
clean:
	docker system prune -f

# Rebuild and restart only backend
restart-backend:
	docker compose build backend && docker compose restart backend
