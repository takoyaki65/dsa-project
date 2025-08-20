import { useMutation, useQuery, useQueryClient, type UseMutationOptions, type UseQueryOptions } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import axios, { type AxiosRequestConfig } from "axios";
import { API_BASE_URL, axiosClient } from "../api/axiosClient";

export interface User {
  id: string;
  name: string;
  email?: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

interface LoginResponse {
  access_token: string;
  token_type: string;
  user: User
}

export const TOKEN_KEY = 'authToken'

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

/**
 * GET request hook with authorization
 * @param queryKey 
 * @param endpoint 
 * @param options 
 * @returns 
 */
export const useAuthQuery = <TData = any>(
  queryKey: string[],
  endpoint: string,
  options?: {
    queryOptions?: Omit<UseQueryOptions<TData, Error>, 'queryKey' | 'queryFn'>;
    axiosConfig?: AxiosRequestConfig;
  }
) => {
  return useQuery<TData, Error>({
    queryKey,
    queryFn: async () => {
      const response = await axiosClient.get<TData>(endpoint, options?.axiosConfig);
      return response.data;
    },
    ...options?.queryOptions,
  });
};

/**
 * POST | PUT | PATCH | DELETE request hook with authorization
 * @param endpoint 
 * @param options 
 * @returns 
 */
export const useAuthMutation = <TData = any, TVariables = any>(
  endpoint: string | ((variables: TVariables) => string),
  options?: {
    method?: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
    mutationOptions?: Omit<UseMutationOptions<TData, Error, TVariables>, 'mutationFn'>;
    axiosConfig?: AxiosRequestConfig;
  }
) => {
  const method = options?.method || 'POST';

  return useMutation<TData, Error, TVariables>({
    mutationFn: async (variables) => {
      const url = typeof endpoint === 'function' ? endpoint(variables) : endpoint;

      let response;
      switch (method) {
        case 'DELETE':
          response = await axiosClient.delete<TData>(url, options?.axiosConfig);
          break;
        case 'PUT':
          response = await axiosClient.put<TData>(url, variables, options?.axiosConfig);
          break;
        case 'PATCH':
          response = await axiosClient.patch<TData>(url, variables, options?.axiosConfig);
          break;
        default:
          response = await axiosClient.post<TData>(url, variables, options?.axiosConfig);
      }

      return response.data;
    },
    ...options?.mutationOptions,
  });
};

/**
 * More flexible fetch hook with authorization
 */
export const useAuthFetch = () => {
  const authFetch = useCallback(async <T = any>(
    config: AxiosRequestConfig
  ): Promise<T> => {
    const response = await axiosClient.request<T>(config);
    return response.data;
  }, []);

  return { authFetch, axiosClient };
};

export const useLogin = () => {
  const { setToken } = useToken();
  const queryClient = useQueryClient();

  return useMutation<LoginResponse, Error, LoginCredentials>({
    mutationFn: async (credentials) => {
      // we request login api without any authorization.
      const response = await axios.post<LoginResponse>(
        `${API_BASE_URL}/login`,
        credentials,
        {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
        }
      );
      return response.data;
    },
    onSuccess: (data) => {
      setToken(data.access_token);
      if (data.user) {
        // save current user information to cache
        queryClient.setQueryData(['currentUser'], data.user);
      }
      // Invalidate all queries to refetch data
      queryClient.invalidateQueries();
    }
  });
};

export const useLogout = () => {
  const { clearToken } = useToken();
  const queryClient = useQueryClient();

  const logoutMutation = useAuthMutation<void, void>(
    '/logout',
    {
      method: 'POST',
      mutationOptions: {
        onSuccess: () => {
          clearToken();
          queryClient.clear();
          window.location.href = '/login'; // Redirect to login page
        },
        onError: () => {
          // clear auth information whatever api fails.
          clearToken();
          queryClient.clear();
          window.location.href = '/login'; // Redirect to login page
        }
      }
    }
  );

  return {
    logout: () => logoutMutation.mutate(),
    isLoading: logoutMutation.isPending,
  };
};

export const useCurrentUser = () => {
  const { token } = useToken();

  return useAuthQuery<User>(
    ['currentUser'],
    'currentUser/me',
    {
      queryOptions: {
        enabled: !!token,
        staleTime: 5 * 60 * 1000, // 5 minutes
        retry: false, // when error encounters, we don't retry
      },
    }
  );
};

export const useAuth = () => {
  const { token } = useToken();
  // const { data: user, isLoading } = useCurrentUser();
  const loginMutation = useLogin();
  const { logout, isLoading: isLogoutLoading } = useLogout();

  return {
    isAuthenticated: !!token,
    login: loginMutation.mutate,
    loginResponse: loginMutation.data,
    logout,
    isLoginPending: loginMutation.isPending,
    isLogoutPending: isLogoutLoading,
    loginError: loginMutation.error,
  }
}
