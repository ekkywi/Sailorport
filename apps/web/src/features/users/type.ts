export type UserRole = "admin" | "developer" | "viewer";

export type User = {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
};