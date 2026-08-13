export type UserRole = "admin" | "developer" | "viewer";

export type User = {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  disabled: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateUserInput = {
  email: string;
  name: string;
  password: string;
  role: UserRole;
};
