from common.database import get_database
from common.service import EmployeeService, UserService


def main():

    database = get_database()

    database.command("ping")

    print("MongoDB connection successful")

    service = EmployeeService()

    employee = service.create({
        "empID": "EMP001",
        "createdBy": "USER123",
        "name": "John Doe",
        "email": "john3@example.com",
        "phone": "+919876543210",
        "department": "IT",
        "salary": 50000,
    })

    print("Created Employee:")
    print(employee)

    employees = service.list(created_by="USER123")

    print("Employees:")
    print(employees)

    user_service = UserService()

    user = user_service.create({
        "name": "Jane Doe",
        "email": "jane3@example.com",
    })

    print("Created User:")
    print(user)

    users = user_service.list()

    print("Users:")
    print(users)

    fetched_user = user_service.get_by_email("jane3@example.com")
    print("Fetched User by Email:")
    print(fetched_user)


if __name__ == "__main__":
    main()
