import { apiClient } from "./employees";
import type { ApiResponse } from "../types/employee";
import type { User } from "../types/auth";

export async function signup(data: any): Promise<any> {
  const res = await apiClient.post<ApiResponse<any>>("/auth/signup", data);
  return res.data;
}

export async function checkStatus(email: string): Promise<string> {
  const res = await apiClient.get<ApiResponse<{ status: string }>>(
    `/auth/status?email=${encodeURIComponent(email)}`
  );
  return res.data.data.status;
}

export async function confirmSignup(data: any): Promise<any> {
  const res = await apiClient.post<ApiResponse<any>>("/auth/confirm", data);
  return res.data;
}

export async function signin(data: any): Promise<any> {
  const res = await apiClient.post<ApiResponse<any>>("/auth/signin", data);
  return res.data;
}

export async function getProfile(): Promise<User> {
  const res = await apiClient.get<ApiResponse<User>>("/profile");
  return res.data.data;
}

export async function logout(): Promise<any> {
  const res = await apiClient.get<ApiResponse<any>>("/auth/logout");
  return res.data;
}

