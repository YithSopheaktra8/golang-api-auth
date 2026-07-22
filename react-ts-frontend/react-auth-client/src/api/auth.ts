import api from "./axios";
import { setAccessToken, clearAccessToken } from "../auth/token";

export interface LoginRequest {
  email: string;

  password: string;
}

export async function login(data: LoginRequest) {
  const response = await api.post("/auth/login", data);

  setAccessToken(response.data.access_token);

  return response.data;
}

export async function logout() {
  await api.post("/auth/logout");

  clearAccessToken();
}
