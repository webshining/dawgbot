-include .env
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))


run-dev:
	python telegram/main.py
rebuild:
	docker compose up -d --no-deps --force-recreate --build ${ARGS}
	