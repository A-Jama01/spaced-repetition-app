export default async function fetchStats() {
  const url = import.meta.env.VITE_API_ADDR + "/v1/stats";
  const response = await fetch(url, {
    method: "GET",
    credentials: "include",
  });
  if (!response.ok) {
    const responseText = await response.text();
    throw new Error(responseText);
  }
  return response.json();
}
