import { QueryClient, QueryClientContext, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";

interface User {
  id: string;
  name: string;
  email?: string;
}

interface LoginCredentials {
  userid: string;
  password: string;
}

interface LoginResponse {
  token: string;
  user: User
}

interface ErrorResponse {
  errors?: {
    [key: string]: any;
  };
}

const TOKEN_KEY = 'authToken'

export const useToken = () => {
  const [token, setTokenState] = useState<string | null>(() => {
    return localStorage.getItem(TOKEN_KEY);
  });

  const setToken = useCallback((newToken: string | null) => {
    if (newToken) {
      localStorage.setItem(TOKEN_KEY, newToken);
    } else {
      localStorage.removeItem(TOKEN_KEY);
    }
    setTokenState(newToken);
  }, []);

  const clearToken = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    setTokenState(null);
  }, []);

  return {
    token,
    setToken,
    clearToken
  };
};

export const useBearer = () => {
  const { token } = useToken();

  const bearerHeader = useCallback((): HeadersInit => {
    if (!token) return {};
    return {
      'Authorization': `Bearer ${token}`
    };
  }, [token]);

  return bearerHeader;
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

/**
 * Custom hook to perform authenticated fetch requests.
 * Automatically includes the Authorization header with the Bearer token.
 * Redirects to login on 401 Unauthorized responses.
 */
export const useAuthFetch = () => {
  const bearerHeader = useBearer();

  const authFetch = useCallback(async <T = any>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> => {
    const url = endpoint.startsWith('http') ? endpoint : `${API_BASE_URL}${endpoint}`;

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...bearerHeader(),
        ...options.headers,
      }
    });

    if (response.status === 401) {
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }

    if (!response.ok) {
      const error: ErrorResponse = await response.json().catch(() => { });
      throw new Error(error.errors?.body || 'Request failed');
    }

    return response.json()
  }, [bearerHeader]);

  return authFetch;
}

export const useLogin = () => {
  const { setToken } = useToken();
  const queryClient = useQueryClient();

  return useMutation<LoginResponse, Error, LoginCredentials>({
    mutationFn: async (credentials) => {
      const response = await fetch(`${API_BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      if (!response.ok) {
        const error: ErrorResponse = await response.json();
        throw new Error(error.errors?.body || "Login failed");
      }

      return response.json();
    },
    onSuccess: (data) => {
      setToken(data.token);

      if (data.user) {
        queryClient.setQueryData(['currentUser'], data.user);
      }
    },
  });
};

export const useLogout = () => {
  const { clearToken } = useToken();
  const queryClient = useQueryClient();

  const logout = useCallback(async () => {
    clearToken();

    queryClient.clear();

    window.location.href = '/login';
  }, [clearToken, queryClient]);

  return { logout };
};

export const useCurrentUser = () => {
  const { token } = useToken();
  const authFetch = useAuthFetch();

  return useQuery<User | null>({
    queryKey: ['currentUser'],
    queryFn: async () => {
      if (!token) return null;
      return authFetch<User>('/user/me');
    },
    enabled: !!token,
    staleTime: 5 * 60 * 1000, // use cache for 5 minutes
  });
};

export const useAuth = () => {
  const { token } = useToken();
  const { data: user, isLoading } = useCurrentUser();
  const loginMutation = useLogin();
  const { logout } = useLogout();

  return {
    isAuthenticated: !!token,
    user,
    isLoading,
    login: loginMutation.mutate,
    logout,
    isLoginPending: loginMutation.isPending,
    loginError: loginMutation.error,
  };
};

export const useApiQuery = <T = any>(
  queryKey: string[],
  endpoint: string,
  options?: {
    enabled?: boolean,
    stableTime?: number,
    refetchInterval?: number,
  }
) => {
  const authFetch = useAuthFetch();

  return useQuery<T>({
    queryKey,
    queryFn: () => authFetch<T>(endpoint),
    ...options,
  });
};

export const useApiMutation = <TData = any, TVariables = any>(
  endpoint: string,
  options?: {
    method?: string;
    onSuccess?: (data: TData) => void;
    onError?: (error: Error) => void
  }
) => {
  const authFetch = useAuthFetch();

  return useMutation<TData, Error, TVariables>({
    mutationFn: (variables) =>
      authFetch<TData>(endpoint, {
        method: options?.method || 'POST',
        body: JSON.stringify(variables),
      }),
    onSuccess: options?.onSuccess,
    onError: options?.onError,
  });
};
