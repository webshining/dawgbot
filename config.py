from environs import Env
from pathlib import Path
import pika


env = Env()
env.read_env()


DIR = Path(__file__).absolute().parent

I18N_PATH = f"{DIR}/locales"
I18N_DOMAIN = env.str("I18N_DOMAIN", "bot")

SURREAL_NS = env.str("SURREAL_NS", "bot")
SURREAL_DB = env.str("SURREAL_DB", "bot")
SURREAL_USER = env.str("SURREAL_USER", None)
SURREAL_PASS = env.str("SURREAL_PASS", None)
SURREAL_URL = env.str("SURREAL_URL", "ws://localhost:8000/rpc")

RABBIT_HOST = env.str("RABBIT_HOST", default="127.0.0.1")
RABBIT_PORT = env.int("RABBIT_PORT", default=5672)
RABBIT_USER = env.str("RABBIT_USER", default=None)
RABBIT_PASS = env.str("RABBIT_PASS", default=None)

RABBIT_CONNECTION_PARAMS = pika.ConnectionParameters(
    host=RABBIT_HOST,
    port=RABBIT_PORT,
    virtual_host="/",
)
if RABBIT_USER and RABBIT_PASS:
    RABBIT_CONNECTION_PARAMS.credentials = pika.PlainCredentials(
        username=RABBIT_USER, password=RABBIT_PASS
    )
