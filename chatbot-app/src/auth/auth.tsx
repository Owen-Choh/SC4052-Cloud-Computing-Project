export interface User {
  userid: number,
	username: string,
}

export interface LoginResponse {
  user: User,
  token: string,
  expiresAt: string,
}
