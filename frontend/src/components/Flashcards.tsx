import type { Flashcard } from "@/pages/Home";
import FlashcardItem from "./FlashcardItem";

interface FlashcardProps {
  cards: Flashcard[];
  deleteCard(cardID: number): Promise<null | Error>;
  updateCard(
    front: string,
    back: string,
    cardID: number,
  ): Promise<null | Error>;
}

export default function Flashcards({
  cards,
  deleteCard,
  updateCard,
}: FlashcardProps) {
  return (
    <div className="grid items-start grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 my-2 mx-4">
      {cards &&
        cards.map((card) => (
          <FlashcardItem
            key={card.id}
            card={card}
            deleteCard={deleteCard}
            updateCard={updateCard}
          />
        ))}
    </div>
  );
}
