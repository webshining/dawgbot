from aiogram import Dispatcher

from .start import router as start_router


def setup_handlers(dp: Dispatcher):
    dp.include_routers(start_router)
