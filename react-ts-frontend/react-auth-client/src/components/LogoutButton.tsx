import { useNavigate } from "react-router-dom";

import { useAuth } from "../auth/AuthContext";

export default function LogoutButton() {
  const auth = useAuth();

  const navigate = useNavigate();

  async function handleLogout() {
    await auth?.logout();

    navigate("/login");
  }

  return <button onClick={handleLogout}>Logout</button>;
}
