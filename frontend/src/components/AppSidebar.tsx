import {
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarFooter,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarInput,
} from "@/components/ui/sidebar";
import { Input } from "./ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import {
  Dialog,
  DialogTrigger,
  DialogHeader,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogClose,
  DialogFooter,
} from "./ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "./ui/button";
import { Label } from "./ui/label";
import { Search, ChartArea, Book, Plus } from "lucide-react";
import { useEffect, useState, type ChangeEvent } from "react";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuTrigger,
  ContextMenuItem,
} from "./ui/context-menu";

export interface Deck {
  id: number;
  user_id: number;
  name: string;
}

export interface SidebarProps {
  decks: Deck[];
  setDecks: React.Dispatch<React.SetStateAction<Deck[]>>;
  fetchDecks(name: String): Promise<void | Error>;
}

export function AppSidebar({ decks, setDecks, fetchDecks }: SidebarProps) {
  const [search, setSearch] = useState<string>("");
  const [deckForm, setDeckForm] = useState<boolean>(false);
  const [deckFormInput, setDeckFormInput] = useState<string>("");
  const [renameInput, setRenameInput] = useState<string>("");

  useEffect(() => {
    fetchDecks(search);
  }, [search]);

  async function createDeck(e: React.FormEvent): Promise<void | Error> {
    e.preventDefault();
    try {
      const url = new URL(import.meta.env.VITE_API_ADDR + "/v1/decks");

      const response = await fetch(url, {
        method: "POST",
        body: JSON.stringify({ name: deckFormInput }),
        credentials: "include",
      });
      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      const data = await response.json();
      setDecks([...decks, data.deck]);
      console.log(data.deck);
      setDeckForm(false);
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error creating deck.");
    }
  }

  async function deleteDeck(id: Number): Promise<void | Error> {
    try {
      const url = new URL(import.meta.env.VITE_API_ADDR + "/v1/decks/" + id);

      const response = await fetch(url, {
        method: "DELETE",
        credentials: "include",
      });
      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      const data = await response.json();
      const filteredDecks = decks.filter((deck) => deck.id !== id);
      setDecks(filteredDecks);
      console.log(data.deck);
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error creating deck.");
    }
  }

  async function renameDeck(id: Number): Promise<void | Error> {
    try {
      const url = new URL(import.meta.env.VITE_API_ADDR + "/v1/decks/" + id);

      const response = await fetch(url, {
        method: "PUT",
        body: JSON.stringify({ name: renameInput }),
        credentials: "include",
      });
      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(responseText);
      }

      const data = await response.json();
      const nextDecks = [...decks];
      const renamedDeck = nextDecks.find((d) => d.id === id);
      if (!renamedDeck) {
        throw new Error("Error renamed deck doesn't exist.");
      }
      renamedDeck!.name = renameInput;
      setDecks(nextDecks);
      setRenameInput("");
    } catch (err) {
      if (err instanceof Error) {
        return err;
      }
      return new Error("Error renaming deck.");
    }
  }

  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarInput
              id="search"
              className="pl-8"
              placeholder="Search for decks..."
              onChange={(e: ChangeEvent<HTMLInputElement>) => {
                setSearch(e.currentTarget.value);
              }}
            />
            <Search className="pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none" />
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton className="cursor-pointer">
              <ChartArea />
              <span>Stats</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Decks</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {decks &&
                decks.map((deck) => (
                  <ContextMenu key={deck.id}>
                    <ContextMenuTrigger asChild>
                      <SidebarMenuButton className="cursor-pointer">
                        <Book />
                        <span>{deck.name}</span>
                      </SidebarMenuButton>
                    </ContextMenuTrigger>
                    <ContextMenuContent>
                      <Dialog>
                        <DialogTrigger asChild>
                          <ContextMenuItem
                            onSelect={(e: Event) => e.preventDefault()}
                          >
                            Rename
                          </ContextMenuItem>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-[425px]">
                          <DialogHeader>
                            <DialogTitle>Rename deck</DialogTitle>
                            <DialogDescription>
                              Save changes to confirm new deck name?
                            </DialogDescription>
                          </DialogHeader>
                          <div className="grid gap-4">
                            <Label htmlFor="username-1">Deck Name</Label>
                            <Input
                              id="deck-rename"
                              name="deck-rename"
                              value={renameInput}
                              onChange={(e: ChangeEvent<HTMLInputElement>) =>
                                setRenameInput(e.currentTarget.value)
                              }
                            />
                          </div>
                          <DialogFooter>
                            <DialogClose asChild>
                              <Button
                                type="submit"
                                onClick={() => renameDeck(deck.id)}
                              >
                                Save changes
                              </Button>
                            </DialogClose>
                            <DialogClose asChild>
                              <Button variant="outline">Cancel</Button>
                            </DialogClose>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                      <ContextMenuItem
                        onSelect={(e: Event) => e.preventDefault()}
                      >
                        <AlertDialog>
                          <AlertDialogTrigger>Delete</AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>
                                Are you absolutely sure you want to delete{" "}
                                <strong>{deck.name}</strong>?
                              </AlertDialogTitle>
                              <AlertDialogDescription>
                                This action cannot be undone. This will
                                permanently delete the deck{" "}
                                <strong>{deck.name}</strong> and all the cards
                                the deck contains.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogAction
                                onClick={() => deleteDeck(deck.id)}
                              >
                                Delete
                              </AlertDialogAction>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </ContextMenuItem>
                    </ContextMenuContent>
                  </ContextMenu>
                ))}
              {deckForm && (
                <form onSubmit={createDeck}>
                  <SidebarMenuItem>
                    <SidebarMenuButton className="cursor-pointer">
                      <Book />
                      <span>
                        <SidebarInput
                          onChange={(e: ChangeEvent<HTMLInputElement>) => {
                            setDeckFormInput(e.currentTarget.value);
                          }}
                        />
                      </span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </form>
              )}
              <SidebarMenuItem>
                <SidebarMenuButton
                  className="cursor-pointer pl-8"
                  onClick={() => {
                    setDeckForm(!deckForm);
                  }}
                >
                  <Plus />
                  <span className="text-zinc-500">New Deck</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
