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