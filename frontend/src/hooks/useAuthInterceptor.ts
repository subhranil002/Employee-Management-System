import { useEffect } from "react";
import type { AuthContextProps } from "react-oidc-context";
import { apiClient } from "../api/employees";

// Attach bearer token and trigger silent token refresh on 401
export function useAuthInterceptor(auth: AuthContextProps) {
  useEffect(() => {
    // Inject access token into outbound requests
    const reqInterceptor = apiClient.interceptors.request.use((config) => {
      if (auth.user?.access_token) {
        config.headers.Authorization = `Bearer ${auth.user.access_token}`;
      }
      return config;
    });

    // Handle token expiry by requesting a silent session renewal
    const resInterceptor = apiClient.interceptors.response.use(
      (res) => res,
      async (err) => {
        if (err.response?.status === 401) {
          if (err.response?.data?.message === "token_expired") {
            try {
              const user = await auth.signinSilent();
              if (user?.access_token) {
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
