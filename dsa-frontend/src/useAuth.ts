import { useState } from "react";
import { getToken, logout } from "./components/LoginPage";

interface UseAuthReturn {
  isAuthenticated: boolean;
  token: string | null;
  logout: () => void;
}

export const useAuth = (): UseAuthReturn => {
  const [token, setToken] = useState<string | null>(getToken());

  return {
    isAuthenticated: !!token,
    token,
    logout: () => {
      logout();
      setToken(null);
    }
  }
}