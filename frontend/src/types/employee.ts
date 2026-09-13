export type Employee = {
  _id: string;
  name: string;
  email: string;
  department: string;
  salary: number;
};

export type EmployeeFormData = Omit<Employee, "_id">;

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};

