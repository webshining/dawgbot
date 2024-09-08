import os
import sys

if os.getenv("ENV", "dev") == "dev":
    sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import asyncio

from loader import dp, bot
from middlewares import setup_middleware
from handlers import setup_handlers


async def main():
    setup_middleware(dp)
    setup_handlers(dp)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
