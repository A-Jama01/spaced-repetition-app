import { AppSidebar } from "@/components/AppSidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { useState } from "react";

export interface Deck {
  id: number;
  user_id: number;
  name: string;
}

export default function Home() {
  const [decks, setDecks] = useState<Deck[]>([]);

  async function fetchDecks(name: string): Promise<void | Error> {
    try {
      const url = new URL(import.meta.env.VITE_API_ADDR + "/v1/decks");
      const params = { name: name };
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
      setDecks(data.decks);
      console.log(data.decks);
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error fetching decks");
    }
  }

  return (
    <SidebarProvider>
      <AppSidebar decks={decks} setDecks={setDecks} fetchDecks={fetchDecks} />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="ml-1" />
          <div>Home</div>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4">
          <div className="grid auto-rows-min gap-4 md:grid-cols-3">
            <div className="bg-muted/50 aspect-video rounded-xl" />
            <div className="bg-muted/50 aspect-video rounded-xl" />
            <div className="bg-muted/50 aspect-video rounded-xl" />
          </div>
          <div className="bg-muted/50 min-h-[100vh] flex-1 rounded-xl md:min-h-min" />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
