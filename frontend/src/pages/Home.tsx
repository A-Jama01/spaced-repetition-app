import { AppSidebar } from "@/components/AppSidebar";
import { Button } from "@/components/ui/button";
import { Card, CardDescription } from "@/components/ui/card";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { useEffect, useState } from "react";

export interface Deck {
  id: number;
  user_id: number;
  name: string;
}

interface Card {
  id: number;
  deckID: number;
  front: string;
  back: string;
}

export default function Home() {
  const [decks, setDecks] = useState<Deck[]>([]);
  const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null);
  const [cards, setCards] = useState<Card[]>([]);
  const [front, setFront] = useState<string>("");
  const [sort, setSort] = useState<string>("");
  const [page, setPage] = useState<number>(1);

  async function fetchCards() {
    if (selectedDeck == null) {
      return;
    }

    try {
      const url = new URL(
        import.meta.env.VITE_API_ADDR +
          "/v1/decks/" +
          selectedDeck.id +
          "/cards",
      );
      const params = {
        front: front,
        sort: sort,
        page: page.toString(),
      };
      url.search = new URLSearchParams(params).toString();

      const response = await fetch(url, {
        method: "GET",
        credentials: "include",
      });
      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }
      const data = await response.json();
      setCards(data.cards);
      console.log(data.cards);
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error fetching cards");
    }
  }
  useEffect(() => {
    fetchCards();
  }, [selectedDeck]);

  return (
    <SidebarProvider>
      <AppSidebar
        decks={decks}
        setDecks={setDecks}
        selectedDeck={selectedDeck}
        setSelectedDeck={setSelectedDeck}
      />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="ml-1" />
          <div>{selectedDeck?.name}</div>
        </header>
        <div className="flex flex-1 flex-col">
          {cards &&
            cards.map((card) => (
              <div>
                <div></div>
              </div>
            ))}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
