from environs import Env
from pathlib import Path


env = Env()
env.read_env()


DIR = Path(__file__).absolute().parent

SURREAL_NS = env.str("SURREAL_NS", "test")
SURREAL_DB = env.str("SURREAL_DB", "test")
SURREAL_USER = env.str("SURREAL_USER", None)
SURREAL_PASS = env.str("SURREAL_PASS", None)
SURREAL_URL = env.str("SURREAL_URL", "ws://localhost:8000/rpc")

RABBIT_HOST = env.str("RABBIT_HOST", default="localhost")
RABBIT_PORT = env.int("RABBIT_PORT", default=5672)
