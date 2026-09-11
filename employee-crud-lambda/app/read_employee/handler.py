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

        employee_id = path_parameters.get(
            "id"
        )

        if employee_id:

            employee = service.get(
                employee_id
            )

            return success(
                200,
                "Employee retrieved successfully",
                employee,
            )

        employees = service.list()

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