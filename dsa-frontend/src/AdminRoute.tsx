import type React from "react";
import { clearStoredToken, useAuth } from "./auth/hooks";
import { Navigate } from "react-router";

interface AdminRouteProps {
  children: React.ReactNode;
}

const AdminRoute: React.FC<AdminRouteProps> = ({ children }) => {
  const { hasAdminScope, isAuthenticated } = useAuth();

  if (!isAuthenticated()) {
    clearStoredToken();
    return <Navigate to="/login" replace />;
  }

  if (!hasAdminScope()) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

export default AdminRoute;
