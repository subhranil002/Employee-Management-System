import { useEffect } from "react";
import type { AuthContextProps } from "react-oidc-context";
import { apiClient } from "../api/employees";

// Custom hook to set up Axios auth headers and handle 401 token refresh/logout
export function useAuthInterceptor(auth: AuthContextProps) {
  useEffect(() => {
    const reqInterceptor = apiClient.interceptors.request.use((config) => {
      if (auth.user?.access_token) {
        config.headers.Authorization = `Bearer ${auth.user.access_token}`;
      }
      return config;
    });

    const resInterceptor = apiClient.interceptors.response.use(
      (res) => res,
      async (err) => {
        if (err.response?.status === 401) {
          if (err.response?.data?.message === "token_expired") {
            try {
              const user = await auth.signinSilent();
              if (user && user.access_token) {
                err.config.headers.Authorization = `Bearer ${user.access_token}`;
                return apiClient.request(err.config);
              }
            } catch {
              auth.removeUser();
            }
          } else {
            auth.removeUser();
          }
        }
        return Promise.reject(err);
      }
    );

    return () => {
      apiClient.interceptors.request.eject(reqInterceptor);
      apiClient.interceptors.response.eject(resInterceptor);
    };
  }, [auth.user?.access_token, auth]);
}

