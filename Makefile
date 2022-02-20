update-static:
	@cd backend/cmd/webapp && rm -rf kodata
	@cd frontend && yarn install && npm run build && mv build kodata && mv kodata ../backend/cmd/webapp/kodata
