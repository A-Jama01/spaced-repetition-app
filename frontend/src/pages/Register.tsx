import { useNavigate } from "react-router";
import { useEffect } from "react";
import RegisterForm from "@/components/RegisterForm";

export default function Register() {
  let navigate = useNavigate();

  const fetchUser = async () => {
    const url = import.meta.env.VITE_API_ADDR + "/v1/auth/me";
    try {
      const response = await fetch(url, {
        method: "GET",
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error(`Response Status: ${response.status}`);
      }

      navigate("/");
    } catch (err) {
      if (err instanceof Error) {
        console.log(err.message);
      } else {
        console.log("Unexpected error", err);
      }
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <div className="grid place-items-center h-screen">
      <RegisterForm onLogin={fetchUser} />
    </div>
  );
}
