export async function fetchUser(): Promise<Error | null> {
  const url = import.meta.env.VITE_API_ADDR + "/v1/auth/me";
  try {
    const response = await fetch(url, {
      method: "GET",
      credentials: "include",
    });
    if (!response.ok) {
      const responseText = await response.text();
      throw new Error(responseText);
    }

    return null;
  } catch (err) {
    if (err instanceof Error) {
      return err;
    }

    return new Error("Error fetching user");
  }
}
