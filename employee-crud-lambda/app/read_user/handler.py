from common.response import error, success
from common.service import (
    UserNotFoundError,
    UserService,
)


service = UserService()


def handler(event, context):

    try:
        query_parameters = (
            event.get("queryStringParameters")
            or {}
        )

        email = query_parameters.get(
            "email"
        )

        if email:

            user = service.get_by_email(
                email
            )

            return success(
                200,
                "User retrieved successfully",
                user,
            )

        users = service.list()

        return success(
            200,
            "Users retrieved successfully",
            users,
        )

    except UserNotFoundError:

        return error(
            404,
            "user not found",
        )

    except Exception as exc:

        print(
            f"Read user error: {exc}"
        )

        return error(
            500,
            "internal server error",
        )
