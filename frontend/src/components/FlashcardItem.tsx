import type { Flashcard } from "@/pages/Home";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Card } from "./ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "./ui/dialog";
import { Button } from "./ui/button";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { Ellipsis, Pencil, Trash2, X } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuGroup,
  DropdownMenuItem,
} from "./ui/dropdown-menu";
import CardForm from "./CardForm";

interface FlashcardItemProps {
  card: Flashcard;
  deleteCard(cardID: number): Promise<null | Error>;
  updateCard(
    front: string,
    back: string,
    cardID: number,
  ): Promise<null | Error>;
}

export default function FlashcardItem({
  card,
  deleteCard,
  updateCard,
}: FlashcardItemProps) {
  const [cardSide, setCardSide] = useState<"front" | "back">("front");
  const [openDropdown, setOpenDropdown] = useState<boolean>(false);
  const [openEdit, setOpenEdit] = useState<boolean>(false);

  const title = "Edit Flashcard";
  const operationName = "Update";

  function handleDelete(cardID: number): void {
    deleteCard(cardID).then(() => setOpenDropdown(false));
  }

  return (
    <Dialog
      onOpenChange={() => {
        setCardSide("front");
        setOpenEdit(false);
      }}
    >
      <DialogTrigger asChild>
        <Card className="drop-shadow-lg rounded-sm hover:bg-zinc-100 focus-visible:ring-2">
          <div className="p-6">
            <div className="prose prose-sm line-clamp-5">
              <Markdown remarkPlugins={[remarkGfm]}>{card.front}</Markdown>
            </div>
          </div>
        </Card>
      </DialogTrigger>
      <DialogContent className="max-h-9/10 overflow-auto">
        <DropdownMenu open={openDropdown} onOpenChange={setOpenDropdown}>
          <DropdownMenuTrigger asChild>
            <Ellipsis className="rounded-md hover:bg-zinc-200" />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuGroup>
              <DropdownMenuItem onSelect={() => setOpenEdit(!openEdit)}>
                {openEdit ? (
                  <>
                    <X /> Close Edit
                  </>
                ) : (
                  <>
                    <Pencil />
                    Edit
                  </>
                )}
              </DropdownMenuItem>
              <DropdownMenuItem
                onSelect={(e) => {
                  e.preventDefault();
                  handleDelete(card.id);
                }}
              >
                <Trash2 />
                Delete
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
        {openEdit ? (
          <CardForm
            title={title}
            operationName={operationName}
            card={card}
            cardOperation={updateCard}
          >
            <DialogClose asChild>
              <Button variant={"secondary"}>Close</Button>
            </DialogClose>
          </CardForm>
        ) : (
          <div>
            <DialogHeader>
              <DialogTitle className="grid grid-cols-2 gap-1">
                <Button
                  variant={"link"}
                  className={cn(
                    cardSide === "front" && "underline underline-offset-4",
                    "hover:bg-zinc-200",
                  )}
                  onClick={() => setCardSide("front")}
                >
                  Front
                </Button>
                <Button
                  variant={"link"}
                  className={cn(
                    cardSide === "back" && "underline underline-offset-4",
                    "hover:bg-zinc-200",
                  )}
                  onClick={() => setCardSide("back")}
                >
                  Back
                </Button>
              </DialogTitle>
              <DialogDescription></DialogDescription>
            </DialogHeader>
            <div className="grid grid-cols-2 gap-1"></div>
            {cardSide === "front" && (
              <div className="prose prose-sm max-h-[80vh] overflow-auto">
                <Markdown remarkPlugins={[remarkGfm]}>{card.front}</Markdown>
              </div>
            )}
            {cardSide === "back" && (
              <div className="prose prose-sm max-h-[80vh] overflow-auto">
                <Markdown remarkPlugins={[remarkGfm]}>{card.back}</Markdown>
              </div>
            )}
            <div className="prose prose-sm max-h-[80vh] overflow-auto"></div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
