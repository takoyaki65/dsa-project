import { useMutation } from "@tanstack/react-query";
import React, { useState, type ChangeEvent, type FormEvent } from "react";
import { Navigate, useNavigate } from "react-router";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

interface LoginCredentials {
  userid: string;
  password: string;
}

interface LoginResponse {
  token: string;
}

interface ErrorResponse {
  errors?: {
    message?: string;
    [key: string]: any;
  };
}

const loginUser = async (credentials: LoginCredentials): Promise<LoginResponse> => {
  const response = await fetch(`${API_BASE_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.errors?.message || 'Login failed');
  }

  return response.json();
};

const saveToken = (token: string): void => {
  localStorage.setItem('authToken', token);
};

export const getToken = (): string | null => {
  return localStorage.getItem('authToken');
};

export const logout = (): void => {
  localStorage.removeItem('authToken');
};

const LoginPage: React.FC = () => {
  const [credentials, setCredentials] = useState<LoginCredentials>({ userid: '', password: '' });

  const navigate = useNavigate();

  const loginMutation = useMutation<LoginResponse, ErrorResponse, LoginCredentials>({
    mutationFn: loginUser,
    onSuccess: (data) => {
      saveToken(data.token);
      navigate('/dashboard');
    },
    onError: (error) => {
      console.error('Login failed:', error);
    },
  })

  const handleSubmit = (e: FormEvent<HTMLFormElement> | React.MouseEvent<HTMLButtonElement>): void => {
    e.preventDefault();
    if (credentials.userid.trim().length === 0 || credentials.password.trim().length === 0) {
      alert('ユーザーIDとパスワードを入力してください');
      return;
    }

    loginMutation.mutate(credentials);
  }

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setCredentials(prev => ({ ...prev, [name]: value }));
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>): void => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleSubmit(e as any);
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4">
      <div className="max-w-md w-full space-y-8">
        <div className="bg-white rounded-2xl shadow-xl border outline-none p-8">
          <div className="textcenter mb-8">
            <h2 className="text-3xl font-bold text-gray-900">ログイン</h2>
            <p className="mt-2 text-sm text-gray-600">アカウントにログインしてください</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="userid" className="block text-sm font-medium text-gray-700 mb-2">
                ユーザーID
              </label>

              <input
                id="userid"
                name="userid"
                type="text"
                required
                value={credentials.userid}
                onChange={handleInputChange}
                placeholder="ユーザーIDを入力してください"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition duration-200"
                autoComplete="username"
                disabled={loginMutation.isPending}
              />
            </div>
            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                パスワード
              </label>
              <input
                id="password"
                name="password"
                type="password"
                required
                value={credentials.password}
                onChange={handleInputChange}
                placeholder="パスワードを入力してください"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparen outline-none transition duration-200"
                autoComplete="current-password"
                onKeyDown={handleKeyDown}
                disabled={loginMutation.isPending}
              />
            </div>

            <button
              type="submit"
              className="w-full py-3 px-4 bg-red-50 border border-red-200 rounded-lg"
              disabled={loginMutation.isPending}
              onClick={handleSubmit}
            >
              {loginMutation.isPending ? "ログイン中..." : "ログイン"}
            </button>
          </form>

          {loginMutation.isError && (
            <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">
                {loginMutation.error?.errors?.message || "ログインに失敗しました。"}
              </p>
            </div>
          )}
        </div>

        <div className="text-center">
          <p className="text-sm text-gray-600">トークンの有効時間: 12時間</p>
        </div>
      </div>
    </div>
  )
}

export default LoginPage;