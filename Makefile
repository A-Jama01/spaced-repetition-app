run/backend:
	cd ./backend;\
		direnv exec . go run ./cmd/api/*

run/frontend/dev:
	cd ./frontend;\
		npm run dev

db/start:
	cd ./backend;\
		docker compose up --build
