export default async function fetchStats(
  forecastTimeRange: string,
  deckName: string,
) {
  const url = new URL(`${import.meta.env.VITE_API_ADDR}/v1/stats`);
  url.searchParams.set("tz", Intl.DateTimeFormat().resolvedOptions().timeZone);
  url.searchParams.set("forecast_length", forecastTimeRange);
  if (deckName !== null) {
    url.searchParams.set("deck_name", deckName);
  }

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
