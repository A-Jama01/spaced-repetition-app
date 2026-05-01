import { AppSidebar } from "@/components/AppSidebar";
import CardView from "@/components/CardView";
import ReviewView from "@/components/ReviewView";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue,
} from "@/components/ui/select";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { toast } from "sonner";
import { fetchUser } from "@/utils/fetchUser";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import StatsView from "@/components/StatsView";

export interface Deck {
  id: number;
  user_id: number;
  name: string;
}

export interface Flashcard {
  id: number;
  deckID: number;
  front: string;
  back: string;
}

export enum SidebarSelection {
  Stats = "Stats",
  Settings = "Settings",
  Deck = "Decks",
}

export default function Home() {
  const [currentSelection, setCurrentSelection] = useState<SidebarSelection>(
    SidebarSelection.Stats,
  );
  const [decks, setDecks] = useState<Deck[]>([]);
  const [selectedDeck, setSelectedDeck] = useState<Deck | null>(null);
  const [showReview, setShowReview] = useState<boolean>(false);
  const [cards, setCards] = useState<Flashcard[]>([]);
  const [dueCards, setDueCards] = useState<Flashcard[]>([]);
  const [front, setFront] = useState<string>("");
  const [sort, setSort] = useState<string>("id");
  const [page, setPage] = useState<number>(1);
  const [statsDeck, setStatsDeck] = useState<string>("");

  let navigate = useNavigate();
  async function checkLoginStatus() {
    let err = await fetchUser();
    if (err) {
      navigate("/login");
    }
    return;
  }

  useEffect(() => {
    checkLoginStatus();
  }, []);

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
      setCards(data.cards ?? []);
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
  }, [selectedDeck, sort, front]);

  async function createCard(
    front: string,
    back: string,
  ): Promise<null | Error> {
    if (selectedDeck == null) {
      return null;
    }

    try {
      const url = new URL(
        import.meta.env.VITE_API_ADDR +
          "/v1/decks/" +
          selectedDeck.id +
          "/cards",
      );

      const response = await fetch(url, {
        method: "POST",
        body: JSON.stringify({ front: front, back: back }),
        credentials: "include",
      });

      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      const data = await response.json();
      const newCard: Flashcard = {
        id: data.card.id,
        deckID: data.card.deck_id,
        front: data.card.front,
        back: data.card.back,
      };
      setCards([newCard, ...cards]);
      toast.success("Card created successfully.");

      return null;
    } catch (err) {
      if (err instanceof Error) {
        console.log(err);
        toast.error("Error creating card.");
        return err;
      }
      return new Error("Error creating card");
    }
  }

  async function deleteCard(cardID: number): Promise<null | Error> {
    if (selectedDeck == null) {
      return null;
    }

    try {
      const url = new URL(
        import.meta.env.VITE_API_ADDR +
          "/v1/decks/" +
          selectedDeck.id +
          "/cards/" +
          cardID,
      );

      const response = await fetch(url, {
        method: "DELETE",
        credentials: "include",
      });

      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      setCards((prev) => prev.filter((card) => card.id !== cardID));
      toast.success("Card deleted successfully.");
      return null;
    } catch (err) {
      if (err instanceof Error) {
        toast.error("Error deleting card.");
        return err;
      }
      return new Error("Error deleting card");
    }
  }

  async function updateCard(
    front: string,
    back: string,
    cardID: number,
  ): Promise<null | Error> {
    if (selectedDeck == null) {
      return null;
    }

    try {
      const url = new URL(
        import.meta.env.VITE_API_ADDR +
          "/v1/decks/" +
          selectedDeck.id +
          "/cards/" +
          cardID,
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

      const data = await response.json();
      const updatedCard: Flashcard = {
        id: data.card.id,
        deckID: data.card.id,
        front: data.card.front,
        back: data.card.back,
      };

      const nextCards: Flashcard[] = cards.map((card) => {
        if (card.id === updatedCard.id) {
          return updatedCard;
        } else {
          return card;
        }
      });
      setCards(nextCards);
      toast.success("Card updated successfully.");

      return null;
    } catch (err) {
      if (err instanceof Error) {
        toast.error("Error updating card.");
        return err;
      }
      return new Error("Error updating card");
    }
  }

  async function getDueCards(): Promise<Error | null> {
    if (selectedDeck == null) {
      return null;
    }

    try {
      const url = new URL(
        import.meta.env.VITE_API_ADDR +
          "/v1/decks/" +
          selectedDeck.id +
          "/cards" +
          "/due",
      );

      const response = await fetch(url, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      const data = await response.json();
      setDueCards(data.cards);
      console.log(data.cards);

      return null;
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error fetching due cards");
    }
  }

  async function reviewCard(
    grade: number,
    cardID: number,
  ): Promise<Error | null> {
    if (selectedDeck == null) {
      return null;
    }
    try {
      const url = new URL(
        `${import.meta.env.VITE_API_ADDR}/v1/decks/${selectedDeck.id}/cards/${cardID}/review`,
      );

      const response = await fetch(url, {
        method: "PATCH",
        body: JSON.stringify({ grade: grade }),
        credentials: "include",
      });

      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      setDueCards(dueCards.slice(1));
      toast.success("Card successfully reviewed");
      return null;
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error reviewing card");
    }
  }

  return (
    <SidebarProvider>
      <AppSidebar
        decks={decks}
        setDecks={setDecks}
        setSelectedDeck={setSelectedDeck}
        setSidebarSelection={setCurrentSelection}
        setShowReview={setShowReview}
      />
      <SidebarInset>
        <header className="flex flex-row justify-between h-16 shrink-0 items-center border-b px-2">
          <div>
            <SidebarTrigger className="ml-1" />
            {currentSelection === SidebarSelection.Deck
              ? selectedDeck?.name
              : currentSelection}
          </div>
          {currentSelection === SidebarSelection.Deck && (
            <div>
              <Button onClick={() => setShowReview(!showReview)}>
                {showReview ? "Close Review" : "Review"}
              </Button>
            </div>
          )}
          {currentSelection === SidebarSelection.Stats && (
            <div>
              <Select value={statsDeck} onValueChange={setStatsDeck}>
                <SelectTrigger
                  className="flex w-40 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[767px]/card:hidden"
                  aria-label="Select a value"
                >
                  <SelectValue placeholder="Select a deck" />
                </SelectTrigger>
                <SelectContent className="rounded-xl">
                  <SelectItem value={null as any} className="rounded-lg">
                    None
                  </SelectItem>
                  {decks &&
                    decks.map((deck) => (
                      <SelectItem
                        key={deck.id}
                        value={`${deck.name}`}
                        className="rounded-lg"
                      >
                        {deck.name}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </header>
        <div className="flex flex-1 flex-col">
          {currentSelection === SidebarSelection.Deck && !showReview && (
            <CardView
              cards={cards}
              sort={sort}
              setSort={setSort}
              front={front}
              setFront={setFront}
              createCard={createCard}
              deleteCard={deleteCard}
              updateCard={updateCard}
            />
          )}
          {currentSelection === SidebarSelection.Deck && showReview && (
            <ReviewView
              dueCards={dueCards}
              getDueCards={getDueCards}
              reviewCard={reviewCard}
            />
          )}
          {currentSelection === SidebarSelection.Stats && (
            <StatsView deck={statsDeck} />
          )}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
