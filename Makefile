update-static:
	@cd backend/cmd/webapp && rm -rf kodata
	@cd frontend && yarn install && yarn react-scripts build && mv build kodata && mv kodata ../backend/cmd/webapp/kodata

run:
	@KO_DATA_PATH=backend/cmd/webapp/kodata/ go run backend/cmd/webapp/main.go

debug:
	@make update-static
	@make run
