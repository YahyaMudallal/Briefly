// src/components/ProtectedRoute.tsx
import { Navigate } from "react-router-dom";
import { useUser } from "../context/UserContext";

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useUser();

  if (loading) return <div>Loading...</div>;
  if (!token) return <Navigate to="/auth" replace />;
  return <>{children}</>;
}