from common.response import error, success
from common.service import EmployeeNotFoundError, EmployeeService

service = EmployeeService()


def handler(event, context):
    try:
        path_parameters = event.get("pathParameters") or {}
        employee_id = path_parameters.get("id")
        if not employee_id:
            return error(400, "employee id is required")

        # Extract tenant owner user id from request headers
        headers = {k.lower(): v for k, v in (event.get("headers") or {}).items()}
        created_by = headers.get("x-user-id")
        if not created_by:
            return error(400, "x-user-id header is required")

        # Delete employee document scoped to user
        service.delete(employee_id, created_by)
        return success(200, "Employee deleted successfully", None)

    except EmployeeNotFoundError:
        return error(404, "employee not found")

    except Exception as exc:
        print(f"Delete employee error: {exc}")
        return error(500, "internal server error")
