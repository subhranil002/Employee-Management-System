import json

from common.response import error, success
from common.service import (
    UserService,
    InvalidUserError,
)


service = UserService()


def handler(event, context):

    try:
        body = event.get("body")

        if not body:
            return error(
                400,
                "request body is required",
            )

        if isinstance(body, str):
            body = json.loads(body)

        user = service.create(body)

        return success(
            201,
            "User created successfully",
            user,
        )

    except json.JSONDecodeError:
        return error(
            400,
            "invalid JSON",
        )

    except InvalidUserError as exc:
        return error(
            400,
            str(exc),
        )

    except Exception as exc:
        print(
            f"Create user error: {exc}"
        )

        return error(
            500,
            "internal server error",
        )
