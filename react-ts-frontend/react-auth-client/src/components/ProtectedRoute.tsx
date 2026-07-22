import { Navigate } from "react-router-dom";

import { useAuth } from "../auth/AuthContext";

interface Props {
  children: React.ReactNode;
}

export default function ProtectedRoute({ children }: Props) {
  const auth = useAuth();

  if (auth?.loading) {
    return <div>Loading...</div>;
  }

  if (!auth?.user) {
    return <Navigate to="/login" replace />;
  }

  return children;
}
