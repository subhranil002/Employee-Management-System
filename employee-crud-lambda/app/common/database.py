from pymongo import MongoClient

from .config import MONGODB_URI, MONGODB_DATABASE


_client = None
_database = None


def get_database():
    global _client
    global _database

    if _database is not None:
        return _database

    _client = MongoClient(
        MONGODB_URI,
        serverSelectionTimeoutMS=5000,
        connectTimeoutMS=5000,
        socketTimeoutMS=5000,
        retryWrites=True,
    )

    # Force a connection check.
    _client.admin.command("ping")

    _database = _client[MONGODB_DATABASE]

    return _database