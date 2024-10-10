up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir="db/migrations" up
