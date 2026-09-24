from pydantic import ValidationError

from .schemas import (
    EmployeeCreate,
    EmployeeUpdate,
    UserCreate,
    CognitoUserCreate,
)
from .repository import EmployeeRepository, UserRepository


class InvalidEmployeeError(Exception):
    pass


class EmployeeNotFoundError(Exception):
    pass


class InvalidUserError(Exception):
    pass


class UserNotFoundError(Exception):
    pass


class EmployeeService:

    def __init__(self):
        self.repository = EmployeeRepository()

    def list(self, created_by):
        return self.repository.find_all(created_by)

    def get(self, employee_id, created_by):
        employee = self.repository.find_by_id(
            employee_id,
            created_by,
        )

        if employee is None:
            raise EmployeeNotFoundError()

        return employee

    def create(self, data):
        try:
            validated = EmployeeCreate(**data)
        except ValidationError as e:
            raise InvalidEmployeeError(str(e))

        existing_emp = self.repository.find_by_emp_id(
            validated.empID,
            validated.createdBy,
        )

        if existing_emp is not None:
            raise InvalidEmployeeError(
                "employee with this empID already exists"
            )

        email = validated.email.lower()

        existing = self.repository.find_by_email(
            email,
            validated.createdBy,
        )

        if existing is not None:
            raise InvalidEmployeeError(
                "employee with this email already exists"
            )

        employee = validated.model_dump()
        employee["email"] = email

        return self.repository.insert(employee)

    def update(self, employee_id, data, created_by):
        if not isinstance(data, dict):
            raise InvalidEmployeeError(
                "request body must be a JSON object"
            )

        try:
            validated = EmployeeUpdate(**data)
        except ValidationError as e:
            raise InvalidEmployeeError(str(e))

        update = validated.model_dump(exclude_unset=True)

        if not update:
            raise InvalidEmployeeError("no fields to update")

        if "email" in update:
            email = update["email"].lower()

            existing = self.repository.find_by_email(
                email,
                created_by,
            )

            if (
                existing is not None
                and existing["_id"] != employee_id
            ):
                raise InvalidEmployeeError(
                    "employee with this email already exists"
                )

            update["email"] = email

        employee = self.repository.update_by_id(
            employee_id,
            update,
            created_by,
        )

        if employee is None:
            raise EmployeeNotFoundError()

        return employee

    def delete(self, employee_id, created_by):
        deleted = self.repository.delete_by_id(
            employee_id,
            created_by,
        )

        if deleted == 0:
            raise EmployeeNotFoundError()


class UserService:

    def __init__(self):
        self.repository = UserRepository()

    def list(self):
        return self.repository.find_all()

    def get(self, user_id):
        user = self.repository.find_by_id(user_id)

        if user is None:
            raise UserNotFoundError()

        return user

    def get_by_email(self, email):
        user = self.repository.find_by_email(email.lower())

        if user is None:
            raise UserNotFoundError()

        return user

    def get_by_cognito_sub(self, sub):
        user = self.repository.find_by_cognito_sub(sub)

        if user is None:
            raise UserNotFoundError()

        return user

    def create(self, data):
        """Legacy/API user creation path."""
        try:
            validated = UserCreate(**data)
        except ValidationError as e:
            raise InvalidUserError(str(e))

        email = validated.email.lower()

        existing = self.repository.find_by_email(email)

        if existing is not None:
            raise InvalidUserError(
                "user with this email already exists"
            )

        user = validated.model_dump()
        user["email"] = email

        return self.repository.insert(user)

    def create_from_cognito(self, event):
        """
        Create or upsert a MongoDB user from a Cognito
        Post Confirmation trigger event.
        """
        try:
            attributes = event["request"]["userAttributes"]
        except (KeyError, TypeError) as exc:
            raise InvalidUserError(
                "Cognito event is missing request.userAttributes"
            ) from exc

        if not isinstance(attributes, dict):
            raise InvalidUserError(
                "Cognito userAttributes must be an object"
            )

        email = attributes.get("email")
        name = attributes.get("name")
        cognito_sub = attributes.get("sub")

        # Fallback if Cognito does not have the standard "name" attribute.
        if not name:
            first_name = (
                attributes.get("given_name") or ""
            ).strip()

            last_name = (
                attributes.get("family_name") or ""
            ).strip()

            name = " ".join(
                part
                for part in (first_name, last_name)
                if part
            )

        try:
            validated = CognitoUserCreate(
                cognitoSub=cognito_sub,
                name=name,
                email=email,
            )
        except ValidationError as exc:
            raise InvalidUserError(str(exc)) from exc

        user = validated.model_dump()
        user["email"] = str(validated.email).lower()

        return self.repository.upsert_by_cognito_sub(user)
