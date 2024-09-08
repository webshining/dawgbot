from config import *


TELEGRAM_BOT_TOKEN = env.str("TELEGRAM_BOT_TOKEN", default=None)

I18N_PATH = f"{DIR}/locales"
I18N_DOMAIN = env.str("I18N_DOMAIN", "bot")
