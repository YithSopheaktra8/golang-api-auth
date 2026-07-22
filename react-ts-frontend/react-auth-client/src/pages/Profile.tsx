import { useEffect, useState } from "react";

import { getProfile } from "../api/user";
import LogoutButton from "../components/LogoutButton";

export default function Profile() {
  const [data, setData] = useState<any>();

  useEffect(() => {
    getProfile().then(setData);
  }, []);

  return (
    <div>
      <h1>Profile</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
      <LogoutButton />
    </div>
  );
}
