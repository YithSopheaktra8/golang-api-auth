import { createContext, useContext, useEffect, useState } from "react";

import api from "../api/axios";

import { setAccessToken, clearAccessToken } from "./token";

import { logout as logoutApi } from "../api/auth";

interface AuthContextType {
  user: boolean;

  loading: boolean;

  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: any) {
  const [user, setUser] = useState(false);

  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function restoreSession() {
      try {
        const response = await api.post("/auth/refresh");

        setAccessToken(response.data.access_token);

        setUser(true);
      } catch (error) {
        setUser(false);
      } finally {
        setLoading(false);
      }
    }

    restoreSession();
  }, []);

  async function logout() {
    try {
      await logoutApi();
    } finally {
      clearAccessToken();

      setUser(false);
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
