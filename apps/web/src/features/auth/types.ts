export type AuthUser = {
  id: string;
  email: string;
  name: string;
  role: string;
  created_at: string;
  updated_at: string;
};

export type LoginResponse = {
  token: string;
  user: AuthUser;
};
