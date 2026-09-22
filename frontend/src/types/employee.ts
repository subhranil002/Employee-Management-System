export type Employee = {
  _id: string;
  empID: string;
  createdBy?: string;
  name: string;
  email: string;
  department: string;
  salary: number;
};

export type EmployeeFormData = {
  empID: string;
  name: string;
  email: string;
  department: string;
  salary: number;
};

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};
