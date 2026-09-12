from common.database import get_database
from common.service import EmployeeService


def main():

    database = get_database()

    database.command("ping")

    print("MongoDB connection successful")

    service = EmployeeService()

    employee = service.create({
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+919876543210",
        "department": "IT",
        "salary": 50000,
    })

    print("Created:")
    print(employee)

    employees = service.list()

    print("Employees:")
    print(employees)


if __name__ == "__main__":
    main()