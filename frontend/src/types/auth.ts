export type User = {
  user_sub: string;
  username: string;
  name?: string;
  email?: string;
};

export type AuthState = {
  user: User | null;
  loading: boolean;
  error: string | null;
};

