export type UserProfile = {
  _id: string;
  name: string;
  email: string;
};

export type AuthState = {
  user: UserProfile | null;
  loading: boolean;
  error: string | null;
};
