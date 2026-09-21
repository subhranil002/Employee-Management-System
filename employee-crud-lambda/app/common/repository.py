from bson import ObjectId
from bson.errors import InvalidId
from pymongo import ReturnDocument

from .config import employee_collection, user_collection
from .database import get_database


class EmployeeRepository:

    def __init__(self):
        database = get_database()
        self.collection = database[
            employee_collection
        ]

    def find_all(self, created_by):
        query = {"createdBy": created_by}
        employees = list(
            self.collection
            .find(query)
            .sort("name", 1)
        )

        return [
            self._serialize(employee)
            for employee in employees
        ]

    def find_by_id(self, employee_id, created_by):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return None

        employee = self.collection.find_one(
            {
                "_id": object_id,
                "createdBy": created_by
            }
        )

        if employee is None:
            return None

        return self._serialize(employee)

    def find_by_email(self, email, created_by=None):
        query = {"email": email}
        if created_by:
            query["createdBy"] = created_by
            
        employee = self.collection.find_one(query)

        if employee is None:
            return None

        return self._serialize(employee)

    def find_by_emp_id(self, emp_id, created_by):
        employee = self.collection.find_one(
            {
                "empID": emp_id,
                "createdBy": created_by
            }
        )

        if employee is None:
            return None

        return self._serialize(employee)

    def insert(self, employee):
        result = self.collection.insert_one(
            employee
        )

        employee["_id"] = result.inserted_id

        return self._serialize(employee)

    def update_by_id(
        self,
        employee_id,
        update,
        created_by,
    ):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return None

        employee = (
            self.collection.find_one_and_update(
                {
                    "_id": object_id,
                    "createdBy": created_by
                },
                {
                    "$set": update
                },
                return_document=ReturnDocument.AFTER,
            )
        )

        if employee is None:
            return None

        return self._serialize(employee)

    def delete_by_id(self, employee_id, created_by):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return 0

        result = self.collection.delete_one(
            {
                "_id": object_id,
                "createdBy": created_by
            }
        )

        return result.deleted_count

    @staticmethod
    def _serialize(employee):
        employee["_id"] = str(
            employee["_id"]
        )

        return employee


class UserRepository:

    def __init__(self):
        database = get_database()
        self.collection = database[
            user_collection
        ]

    def find_all(self):
        users = list(
            self.collection
            .find({})
            .sort("name", 1)
        )

        return [
            self._serialize(user)
            for user in users
        ]

    def find_by_id(self, user_id):
        try:
            object_id = ObjectId(user_id)
        except InvalidId:
            return None

        user = self.collection.find_one(
            {
                "_id": object_id
            }
        )

        if user is None:
            return None

        return self._serialize(user)

    def find_by_email(self, email):
        user = self.collection.find_one(
            {
                "email": email
            }
        )

        if user is None:
            return None

        return self._serialize(user)

    def insert(self, user):
        result = self.collection.insert_one(
            user
        )

        user["_id"] = str(result.inserted_id)

        return user

    @staticmethod
    def _serialize(user):
        user["_id"] = str(
            user["_id"]
        )

        return user
