import type { Flashcard } from "@/pages/Home";
import { Card, CardContent, CardFooter } from "./ui/card";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Button } from "./ui/button";
import { useEffect, useState } from "react";

interface ReviewViewProps {
  dueCards: Flashcard[];
  getDueCards(): Promise<Error | null>;
  reviewCard(grade: number, cardID: number): Promise<Error | null>;
}

export default function ReviewView({
  dueCards,
  getDueCards,
  reviewCard,
}: ReviewViewProps) {
  const [showBack, setShowBack] = useState<boolean>(false);

  function handleReviewClick(e: React.MouseEvent<HTMLButtonElement>): void {
    //TODO: send grade to server then show next card that is due
    if (dueCards) {
      const grade = Number(e.currentTarget.value);
      const cardID = dueCards[0].id;

      reviewCard(grade, cardID);
      setShowBack(false);
    }
  }

  useEffect(() => {
    getDueCards();
  }, []);

  if (dueCards === undefined || dueCards === null || dueCards.length === 0) {
    return "No cards due";
  }

  return (
    <div className="flex flex-col min-h-full items-center">
      <Card className="flex flex-col w-full max-w-xl max-h-[70vh] mt-10">
        <CardContent className="flex-1 overflow-y-auto">
          {!showBack ? (
            <div className="prose prose-sm">
              <Markdown remarkPlugins={[remarkGfm]}>
                {dueCards[0].front}
              </Markdown>
            </div>
          ) : (
            <div className="prose prose-sm">
              <Markdown remarkPlugins={[remarkGfm]}>
                {dueCards[0].back}
              </Markdown>
            </div>
          )}
        </CardContent>
        <CardFooter className="flex justify-center">
          <Button variant="outline" onClick={() => setShowBack(!showBack)}>
            {showBack ? "Show Front" : "Show Back"}
          </Button>
        </CardFooter>
      </Card>
      {showBack && (
        <div className="mt-auto flex gap-6 mb-6">
          <Button value={1} onClick={handleReviewClick}>
            Forgot
          </Button>
          <Button value={2} onClick={handleReviewClick}>
            Hard
          </Button>
          <Button value={3} onClick={handleReviewClick}>
            Good
          </Button>
          <Button value={4} onClick={handleReviewClick}>
            Easy
          </Button>
        </div>
      )}
    </div>
  );
}
