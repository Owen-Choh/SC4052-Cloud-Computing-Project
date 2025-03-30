export interface User {
  userid: number;
  username: string;
}

export interface LoginResponse {
  user: User;
  expiresAt: string;
}
