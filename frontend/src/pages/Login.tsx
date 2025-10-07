import { useEffect } from "react";
import LoginForm from "../components/LoginForm";
import { useNavigate } from "react-router";
import { fetchUser } from "@/utils/fetchUser";

export default function Login() {
  let navigate = useNavigate();

  async function handleValidUser() {
    let err = await fetchUser();
    if (err === null) {
      navigate("/");
    }
  }

  useEffect(() => {
    handleValidUser();
  }, []);

  return (
    <div className="grid place-items-center h-screen">
      <LoginForm onLogin={handleValidUser} />
    </div>
  );
}
