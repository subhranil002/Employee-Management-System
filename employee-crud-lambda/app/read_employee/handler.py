from common.response import error, success
from common.service import (
    EmployeeNotFoundError,
    EmployeeService,
)


service = EmployeeService()


def handler(event, context):

    try:
        path_parameters = (
            event.get("pathParameters")
            or {}
        )

        headers = {k.lower(): v for k, v in (event.get("headers") or {}).items()}
        created_by = headers.get("x-user-id")

        if not created_by:
            return error(
                400,
                "x-user-id header is required",
            )

        employee_id = path_parameters.get(
            "id"
        )

        if employee_id:

            employee = service.get(
                employee_id,
                created_by
            )

            return success(
                200,
                "Employee retrieved successfully",
                employee,
            )

        employees = service.list(created_by)

        return success(
            200,
            "Employees retrieved successfully",
            employees,
        )

    except EmployeeNotFoundError:

        return error(
            404,
            "employee not found",
        )

    except Exception as exc:

        print(
            f"Read employee error: {exc}"
        )

        return error(
            500,
            "internal server error",
        )
