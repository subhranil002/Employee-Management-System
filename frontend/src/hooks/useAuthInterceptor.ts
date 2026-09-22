import { useEffect } from "react";
import type { AuthContextProps } from "react-oidc-context";
import { apiClient } from "../api/employees";

// Attaches Authorization + X-Id-Token headers and handles 401 silent refresh
export function useAuthInterceptor(auth: AuthContextProps) {
  useEffect(() => {
    const reqInterceptor = apiClient.interceptors.request.use((config) => {
      if (auth.user?.access_token) {
        config.headers.Authorization = `Bearer ${auth.user.access_token}`;
      }
      if (auth.user?.id_token) {
        config.headers["X-Id-Token"] = auth.user.id_token;
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
              if (user?.access_token) {
                err.config.headers.Authorization = `Bearer ${user.access_token}`;
                if (user.id_token) {
                  err.config.headers["X-Id-Token"] = user.id_token;
                }
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
  }, [auth.user?.access_token, auth.user?.id_token, auth]);
}
