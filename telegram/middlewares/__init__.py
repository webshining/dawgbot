from aiogram import Dispatcher

from .user import UserMiddleware
from .i18n import i18n_middleware


def setup_middleware(dp: Dispatcher):
    dp.update.middleware(i18n_middleware)
    dp.update.middleware(UserMiddleware())
