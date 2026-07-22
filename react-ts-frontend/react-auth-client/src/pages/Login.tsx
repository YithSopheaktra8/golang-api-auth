import { useState } from "react";

import { login } from "../api/auth";

import { useNavigate } from "react-router-dom";

export default function Login() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");

  const [password, setPassword] = useState("");

  async function submit() {
    await login({
      email,

      password,
    });

    navigate("/");
  }

  return (
    <div>
      <h1>Login</h1>

      <input placeholder="email" onChange={(e) => setEmail(e.target.value)} />

      <input
        placeholder="password"
        type="password"
        onChange={(e) => setPassword(e.target.value)}
      />

      <button onClick={submit}>Login</button>
    </div>
  );
}
