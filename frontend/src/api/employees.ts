import axios from "axios";
import toast from "react-hot-toast";
import type {
  ApiResponse,
  Employee,
  EmployeeFormData,
} from "../types/employee";

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

export const apiClient = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

export async function getEmployees(): Promise<Employee[]> {
  const res = await apiClient.get<ApiResponse<Employee[]>>("/employees");
  return res.data.data;
}

export async function getEmployee(id: string): Promise<Employee> {
  const res = await apiClient.get<ApiResponse<Employee>>(`/employees/${id}`);
  return res.data.data;
}

export async function createEmployee(
  data: EmployeeFormData,
): Promise<Employee> {
  const res = apiClient.post<ApiResponse<Employee>>("/employees", data);
  toast.promise(res, {
    loading: "Saving employee...",
    success: (r) => r.data.message,
    error: (err) => err?.response?.data?.message || "Failed to save",
  });
  return (await res).data.data;
}

export async function updateEmployee(
  id: string,
  data: EmployeeFormData,
): Promise<Employee> {
  const res = apiClient.patch<ApiResponse<Employee>>(`/employees/${id}`, data);
  toast.promise(res, {
    loading: "Updating employee...",
    success: (r) => r.data.message,
    error: (err) => err?.response?.data?.message || "Failed to update",
  });
  return (await res).data.data;
}

export async function deleteEmployee(id: string): Promise<null> {
  const res = apiClient.delete<ApiResponse<null>>(`/employees/${id}`);
  toast.promise(res, {
    loading: "Deleting employee...",
    success: (r) => r.data.message,
    error: (err) => err?.response?.data?.message || "Failed to delete",
  });
  return (await res).data.data;
}

export async function getHealth(): Promise<null> {
  const res = await apiClient.get<ApiResponse<null>>("/health");
  return res.data.data;
}
