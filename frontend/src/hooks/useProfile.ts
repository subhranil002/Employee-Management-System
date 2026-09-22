import { useEffect, useState } from "react";
import type { UserProfile } from "../types/auth";
import { getProfile } from "../api/user";

type ProfileState = {
  profile: UserProfile | null;
  loading: boolean;
  error: string;
};

// Fetches the authenticated user profile from the backend /profile endpoint
export function useProfile(isAuthenticated: boolean) {
  const [state, setState] = useState<ProfileState>({
    profile: null,
    loading: true,
    error: "",
  });

  useEffect(() => {
    let ignore = false;
    if (!isAuthenticated) {
      setState({ profile: null, loading: false, error: "" });
      return;
    }

    setState((s) => ({ ...s, loading: true, error: "" }));

    getProfile()
      .then((data) => {
        if (!ignore) setState({ profile: data, loading: false, error: "" });
      })
      .catch((err) => {
        if (!ignore) {
          const msg = err?.response?.data?.message || err.message || "Failed to load profile";
          setState({ profile: null, loading: false, error: msg });
        }
      });

    return () => {
      ignore = true;
    };
  }, [isAuthenticated]);

  return state;
}

