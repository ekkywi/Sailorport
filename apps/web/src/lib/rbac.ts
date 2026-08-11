export function isAdmin(role: string): boolean {
  return role === "admin";
}

export function canWriteCatalog(role: string): boolean {
  return role === "admin" || role === "developer";
}