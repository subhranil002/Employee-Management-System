from common.response import error, success
from common.service import UserNotFoundError, UserService

service = UserService()


def handler(event, context):
    try:
        query_parameters = event.get("queryStringParameters") or {}
        email = query_parameters.get("email")
        sub = query_parameters.get("sub")

        # Query single user when sub or email query parameter is provided
        if sub:
            user = service.get_by_cognito_sub(sub)
            return success(200, "User retrieved successfully", user)
        elif email:
            user = service.get_by_email(email)
            return success(200, "User retrieved successfully", user)

        # Otherwise retrieve full user listing
        users = service.list()
        return success(200, "Users retrieved successfully", users)

    except UserNotFoundError:
        return error(404, "user not found")

    except Exception as exc:
        print(f"Read user error: {exc}")
        return error(500, "internal server error")
