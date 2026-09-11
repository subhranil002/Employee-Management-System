from bson import ObjectId
from bson.errors import InvalidId
from pymongo import ReturnDocument

from .config import MONGODB_COLLECTION
from .database import get_database


class EmployeeRepository:

    def __init__(self):
        database = get_database()
        self.collection = database[
            MONGODB_COLLECTION
        ]

    def find_all(self):
        employees = list(
            self.collection
            .find({})
            .sort("name", 1)
        )

        return [
            self._serialize(employee)
            for employee in employees
        ]

    def find_by_id(self, employee_id):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return None

        employee = self.collection.find_one(
            {
                "_id": object_id
            }
        )

        if employee is None:
            return None

        return self._serialize(employee)

    def find_by_email(self, email):
        employee = self.collection.find_one(
            {
                "email": email
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
    ):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return None

        employee = (
            self.collection.find_one_and_update(
                {
                    "_id": object_id
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

    def delete_by_id(self, employee_id):
        try:
            object_id = ObjectId(employee_id)
        except InvalidId:
            return 0

        result = self.collection.delete_one(
            {
                "_id": object_id
            }
        )

        return result.deleted_count

    @staticmethod
    def _serialize(employee):
        employee["_id"] = str(
            employee["_id"]
        )

        return employee