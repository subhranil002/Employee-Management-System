from common.response import error, success
from common.service import (
    UserNotFoundError,
    UserService,
)


service = UserService()


def handler(event, context):

    try:
        path_parameters = (
            event.get("pathParameters")
            or {}
        )

        user_id = path_parameters.get(
            "id"
        )

        if user_id:

            user = service.get(
                user_id
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
