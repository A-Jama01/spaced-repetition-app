import { ArrowDownAz, ChevronDown, Plus, Search } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import Flashcards from "./Flashcards";
import { type Flashcard } from "@/pages/Home";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { DropdownMenuRadioItem } from "./ui/dropdown-menu";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogTitle,
  DialogTrigger,
  DialogHeader,
  DialogDescription,
} from "./ui/dialog";
import { Textarea } from "./ui/textarea";
import { useState } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {
  Field,
  FieldGroup,
  FieldSet,
  FieldLabel,
  FieldDescription,
  FieldTitle,
} from "./ui/field";
import CardForm from "./CardForm";

interface CardViewProps {
  cards: Flashcard[];
  deckID: number;
  sort: string;
  setSort: React.Dispatch<React.SetStateAction<string>>;
  front: string;
  setFront: React.Dispatch<React.SetStateAction<string>>;
  createCard(front: string, back: string): Promise<void>;
  deleteCard(cardID: number): Promise<null | Error>;
  updateCard(
    front: string,
    back: string,
    cardID: number,
  ): Promise<null | Error>;
}

export default function CardView({
  cards,
  deckID,
  sort,
  setSort,
  front,
  setFront,
  createCard,
  deleteCard,
  updateCard,
}: CardViewProps) {
  const title = "Create a FlashCard";
  const operationName = "Create";

  return (
    <div className="flex flex-col h-full">
      <div className="flex justify-between my-2 mx-3">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant={"outline"} size={"sm"}>
              <ArrowDownAz />
              Sort
              <ChevronDown />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuRadioGroup value={sort} onValueChange={setSort}>
              <DropdownMenuRadioItem value="id">None</DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="-created_at">
                Newest
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="created_at">
                Oldest
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="last_review">
                Most Recently Reviewed
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="-last_review">
                Least Recently Reviewed
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="due">
                Earliest Due
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="-due">
                Latest Due
              </DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
        <div className="relative w-full max-w-sm">
          <Search className="absolute top-2 left-1 size-5 opacity-50" />
          <Input
            className="pl-8"
            type="search"
            placeholder="Search for cards..."
            value={front}
            onChange={(e) => setFront(e.target.value)}
          />
        </div>
        <Dialog>
          <DialogTrigger asChild>
            <Button variant={"outline"} size={"sm"}>
              New Card
              <Plus />
            </Button>
          </DialogTrigger>
          <DialogContent className="max-h-9/10 overflow-auto">
            <DialogHeader>
              <DialogTitle></DialogTitle>
              <DialogDescription></DialogDescription>
            </DialogHeader>
            <CardForm
              title={title}
              operationName={operationName}
              cardOperation={createCard}
            >
              <DialogClose asChild>
                <Button variant={"secondary"}>Close</Button>
              </DialogClose>
            </CardForm>
          </DialogContent>
        </Dialog>
      </div>
      <div className="flex-1 overflow-y-auto">
        <Flashcards
          cards={cards}
          deckID={deckID}
          front={front}
          sort={sort}
          deleteCard={deleteCard}
          updateCard={updateCard}
        />
      </div>
    </div>
  );
}
