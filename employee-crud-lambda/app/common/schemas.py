from typing import Optional

from pydantic import BaseModel, EmailStr, Field


class EmployeeCreate(BaseModel):
    empID: str = Field(
        ...,
        min_length=1,
        description="Employee ID",
    )
    createdBy: str = Field(
        ...,
        min_length=1,
        description="ID of the user who created this employee",
    )
    name: str = Field(
        ...,
        min_length=1,
        description="Full name of the employee",
    )
    email: EmailStr
    department: str = Field(
        ...,
        min_length=1,
        description="Department name",
    )
    salary: float = Field(
        ...,
        ge=0.0,
        description="Salary must be non-negative",
    )


class EmployeeUpdate(BaseModel):
    name: Optional[str] = Field(
        None,
        min_length=1,
    )
    email: Optional[EmailStr] = None
    department: Optional[str] = Field(
        None,
        min_length=1,
    )
    salary: Optional[float] = Field(
        None,
        ge=0.0,
    )


class UserCreate(BaseModel):
    name: str = Field(
        ...,
        min_length=1,
        description="User's name",
    )
    email: EmailStr


class CognitoUserCreate(BaseModel):
    cognitoSub: str = Field(
        ...,
        min_length=1,
        description="Cognito user sub (unique identifier)",
    )
    name: str = Field(
        ...,
        min_length=1,
        description="User's full name",
    )
    email: EmailStr
