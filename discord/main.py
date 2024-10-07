import os
import sys

if os.getenv("ENV", "dev") == "dev":
    sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from handlers import bot
from dsconfig import DISCORD_BOT_TOKEN


if __name__ == "__main__":
    bot.run(DISCORD_BOT_TOKEN)
