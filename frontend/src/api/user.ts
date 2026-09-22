import type { ApiResponse } from "../types/employee";
import type { UserProfile } from "../types/auth";
import { apiClient } from "./employees";

// Fetches the resolved user profile from the backend (verifies tokens + get-or-create)
export async function getProfile(): Promise<UserProfile> {
  const res = await apiClient.get<ApiResponse<UserProfile>>("/profile");
  return res.data.data;
}

