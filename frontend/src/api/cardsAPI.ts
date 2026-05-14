export async function fetchCards(
  deckID: number,
  front: string,
  sort: string,
  page: number,
): Promise<any> {
  const url = new URL(
    `${import.meta.env.VITE_API_ADDR}/v1/decks/${deckID}/cards`,
  );
  url.searchParams.set("front", front);
  url.searchParams.set("sort", sort);
  url.searchParams.set("page", page.toString());

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

export async function deleteCard(
  deckID: number,
  cardID: number,
): Promise<void> {
  const url = new URL(
    `${import.meta.env.VITE_API_ADDR}/v1/decks/${deckID}/cards/${cardID}`,
  );

  const response = await fetch(url, {
    method: "DELETE",
    credentials: "include",
  });
  if (!response.ok) {
    const responseText = await response.text();
    throw new Error(responseText);
  }
}

export async function updateCard(
  deckID: number,
  cardID: number,
  front: string,
  back: string,
): Promise<any> {
  const url = new URL(
    `${import.meta.env.VITE_API_ADDR}/v1/decks/${deckID}/cards/${cardID}`,
  );

  const response = await fetch(url, {
    method: "PATCH",
    body: JSON.stringify({ front: front, back: back }),
    credentials: "include",
  });

  if (!response.ok) {
    const responseText = await response.text();
    throw new Error(responseText);
  }

  return response.json();
}
