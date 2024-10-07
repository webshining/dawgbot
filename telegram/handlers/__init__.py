from aiogram import Dispatcher

from .start import router as start_router
from .notify import router as auth_router
from .notifications import router as notifications_router
from .remove import router as remove_router


def setup_handlers(dp: Dispatcher):
    dp.include_routers(auth_router, start_router, notifications_router, remove_router)
