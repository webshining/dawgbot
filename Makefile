-include .env
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
LOCALES_PATH := ./locales
I18N_DOMAIN ?= bot

.PHONY: discord telegram

dev-telegram:
	python telegram/main.py
dev-discord:
	python discord/main.py
rebuild:
	docker compose up -d --no-deps --force-recreate --build ${ARGS}
pybabel_extract:
	pybabel extract --input-dirs=. -o $(LOCALES_PATH)/$(I18N_DOMAIN).pot
pybabel_init: 
	pybabel init -i $(LOCALES_PATH)/$(I18N_DOMAIN).pot -d $(LOCALES_PATH) -D $(I18N_DOMAIN) -l en && \
	pybabel init -i $(LOCALES_PATH)/$(I18N_DOMAIN).pot -d $(LOCALES_PATH) -D $(I18N_DOMAIN) -l ru && \
	pybabel init -i $(LOCALES_PATH)/$(I18N_DOMAIN).pot -d $(LOCALES_PATH) -D $(I18N_DOMAIN) -l uk
pybabel_update: 
	pybabel update -i $(LOCALES_PATH)/$(I18N_DOMAIN).pot -d $(LOCALES_PATH) -D $(I18N_DOMAIN)
pybabel_compile:
	pybabel compile -d $(LOCALES_PATH) -D $(I18N_DOMAIN)
db_export:
	docker compose exec surrealdb //surreal export --conn http://localhost:8000 --user $(SURREAL_USER) --pass $(SURREAL_PASS) --ns $(SURREAL_NS) --db $(SURREAL_DB) export.surql && \
	docker compose cp surrealdb:/export.surql ./export.surql
db_import:
	docker compose cp ./export.surql surrealdb:/export.surql && \
	docker compose exec surrealdb /surreal import --conn http://localhost:8000 --user $(SURREAL_USER) --pass $(SURREAL_PASS) --ns $(SURREAL_NS) --db $(SURREAL_DB) export.surql