import { Link } from "react-router";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  CardAction,
} from "./ui/card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useState } from "react";

interface RegisterProps {
  onRegister: () => Promise<void>;
}

export default function RegisterForm({ onRegister }: RegisterProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function register(e: React.FormEvent) {
    e.preventDefault();

    const url = import.meta.env.VITE_API_ADDR + "/v1/auth/register";
    try {
      const response = await fetch(url, {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({ username: username, password: password }),
      });
      if (!response.ok) {
        throw new Error(`Response Status: ${response.status}`);
      }

      const result = await response.json();
      console.log(result);

      onRegister();
    } catch (err) {
      if (err instanceof Error) {
        console.log(err.message);
      } else {
        console.log("Unexpected error", err);
      }
    }
  }

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <div className="grid grid-cols-1 gap-4">
          <CardTitle>Register an account</CardTitle>
          <CardDescription>
            Enter a username and password to register an account.
          </CardDescription>
        </div>
        <div>
          <CardAction>
            <Link to="/login">
              <Button variant="link">Login</Button>
            </Link>
          </CardAction>
        </div>
      </CardHeader>
      <form onSubmit={register}>
        <CardContent>
          <div className="flex flex-col gap-6">
            <div className="grid gap-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                name="username"
                type="text"
                placeholder="Username"
                required
                minLength={8}
                maxLength={40}
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                name="password"
                type="text"
                placeholder="Password"
                required
                minLength={8}
                maxLength={50}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
            </div>
          </div>
        </CardContent>
        <CardFooter>
          <Button type="submit" className="w-full">
            Register
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
