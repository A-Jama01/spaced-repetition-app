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
import { Alert, AlertTitle, AlertDescription } from "./ui/alert";
import { AlertCircleIcon } from "lucide-react";

interface LoginProps {
  onLogin: () => void;
}

export default function LoginForm({ onLogin }: LoginProps) {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [showAlert, setShowAlert] = useState<boolean>(false);

  async function login(e: React.FormEvent) {
    e.preventDefault();

    const url = import.meta.env.VITE_API_ADDR + "/v1/auth/login";
    try {
      const response = await fetch(url, {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({ username: username, password: password }),
      });
      if (!response.ok) {
        setShowAlert(true);
        throw new Error(`Response Status: ${response.status}`);
      }

      const result = await response.json();
      console.log(result);
      onLogin();
    } catch (err) {
      if (err instanceof Error) {
        console.log(err.message);
      } else {
        console.log("Unexpected error", err);
      }
    }
  }

  return (
    <div>
      <Card className="w-full max-w-sm">
        <CardHeader>
          <div className="grid grid-cols-1 gap-4">
            <CardTitle>Login to your account</CardTitle>
            <CardDescription>
              Enter username and password to login to your account.
            </CardDescription>
          </div>
          <div>
            <CardAction>
              <Link to="/register">
                <Button variant="link">Sign Up</Button>
              </Link>
            </CardAction>
          </div>
        </CardHeader>
        <form onSubmit={login}>
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
              Login
            </Button>
          </CardFooter>
        </form>
      </Card>
      {showAlert && (
        <Alert className="mt-4 w-full max-w-sm" variant="destructive">
          <AlertCircleIcon />
          <AlertTitle>Unable to login to account</AlertTitle>
          <AlertDescription>
            <p>Please verify that your login information is correct.</p>
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
