# Makefile

# Start backend using reflex (for live reload)
run-be:
	cd backend && go run main.go

# Start frontend (React/Vite)
run-fe:
	cd frontend/client && npm start
