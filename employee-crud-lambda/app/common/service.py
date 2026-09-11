import re


from .repository import EmployeeRepository


EMAIL_REGEX = re.compile(
    r"^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$",
    re.IGNORECASE,
)


class InvalidEmployeeError(Exception):
    pass


class EmployeeNotFoundError(Exception):
    pass


class EmployeeService:

    def __init__(self):
        self.repository = EmployeeRepository()

    def list(self):
        return self.repository.find_all()

    def get(self, employee_id):
        employee = (
            self.repository.find_by_id(
                employee_id
            )
        )

        if employee is None:
            raise EmployeeNotFoundError()

        return employee

    def create(self, data):

        self._validate_name(
            data.get("name")
        )

        self._validate_email(
            data.get("email")
        )

        self._validate_department(
            data.get("department")
        )

        self._validate_salary(
            data.get("salary")
        )

        name = data["name"].strip()

        email = (
            data["email"]
            .strip()
            .lower()
        )

        department = (
            data["department"]
            .strip()
        )

        salary = float(
            data["salary"]
        )

        existing = (
            self.repository.find_by_email(
                email
            )
        )

        if existing is not None:
            raise InvalidEmployeeError(
                "employee with this email already exists"
            )

        employee = {
            "name": name,
            "email": email,
            "department": department,
            "salary": salary,
        }

        return self.repository.insert(
            employee
        )

    def update(
        self,
        employee_id,
        data,
    ):

        if not isinstance(data, dict):
            raise InvalidEmployeeError(
                "request body must be a JSON object"
            )

        update = {}

        if "name" in data:
            self._validate_name(
                data["name"]
            )

            update["name"] = (
                data["name"].strip()
            )

        if "email" in data:
            self._validate_email(
                data["email"]
            )

            email = (
                data["email"]
                .strip()
                .lower()
            )

            existing = (
                self.repository.find_by_email(
                    email
                )
            )

            if (
                existing is not None
                and existing["_id"] != employee_id
            ):
                raise InvalidEmployeeError(
                    "employee with this email already exists"
                )

            update["email"] = email

        if "department" in data:
            self._validate_department(
                data["department"]
            )

            update["department"] = (
                data["department"].strip()
            )

        if "salary" in data:
            self._validate_salary(
                data["salary"]
            )

            update["salary"] = float(
                data["salary"]
            )

        if not update:
            raise InvalidEmployeeError(
                "no fields to update"
            )

        employee = (
            self.repository.update_by_id(
                employee_id,
                update,
            )
        )

        if employee is None:
            raise EmployeeNotFoundError()

        return employee

    def delete(self, employee_id):

        deleted = (
            self.repository.delete_by_id(
                employee_id
            )
        )

        if deleted == 0:
            raise EmployeeNotFoundError()

    @staticmethod
    def _validate_name(name):

        if not isinstance(name, str):
            raise InvalidEmployeeError(
                "name is required"
            )

        if not name.strip():
            raise InvalidEmployeeError(
                "name is required"
            )

    @staticmethod
    def _validate_email(email):

        if not isinstance(email, str):
            raise InvalidEmployeeError(
                "email is required"
            )

        email = email.strip().lower()

        if not email:
            raise InvalidEmployeeError(
                "email is required"
            )

        if not EMAIL_REGEX.match(email):
            raise InvalidEmployeeError(
                "invalid email format"
            )

    @staticmethod
    def _validate_department(department):

        if not isinstance(
            department,
            str,
        ):
            raise InvalidEmployeeError(
                "department is required"
            )

        if not department.strip():
            raise InvalidEmployeeError(
                "department is required"
            )

    @staticmethod
    def _validate_salary(salary):

        if salary is None:
            raise InvalidEmployeeError(
                "salary is required"
            )

        try:
            salary = float(salary)
        except (TypeError, ValueError):
            raise InvalidEmployeeError(
                "invalid salary"
            )

        if salary < 0:
            raise InvalidEmployeeError(
                "salary cannot be negative"
            )