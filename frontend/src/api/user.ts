import type { ApiResponse } from "../types/employee";
import type { UserProfile } from "../types/auth";
import { apiClient } from "./employees";

// Fetch authenticated user profile from backend
export async function getProfile(): Promise<UserProfile> {
  const res = await apiClient.get<ApiResponse<UserProfile>>("/profile");
  return res.data.data;
}
