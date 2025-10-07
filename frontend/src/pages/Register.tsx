import { useNavigate } from "react-router";
import { fetchUser } from "@/utils/fetchUser";
import { useEffect } from "react";
import RegisterForm from "@/components/RegisterForm";

export default function Register() {
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
      <RegisterForm onRegister={handleValidUser} />
    </div>
  );
}
