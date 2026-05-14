import type { Flashcard } from "@/pages/Home";
import FlashcardItem from "./FlashcardItem";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { fetchCards, deleteCard, updateCard } from "@/api/cardsAPI";
import { toast } from "sonner";

interface FlashcardProps {
  cards: Flashcard[];
  deckID: number;
  sort: string;
  front: string;
  deleteCard(cardID: number): Promise<null | Error>;
  updateCard(
    front: string,
    back: string,
    cardID: number,
  ): Promise<null | Error>;
}

export default function Flashcards({
  cards,
  deckID,
  front,
  sort,
  // deleteCard,
  // updateCard,
}: FlashcardProps) {
  const queryClient = useQueryClient();
  const queryKey = ["cards", deckID, front, sort];
  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
    status,
  } = useInfiniteQuery({
    queryKey: queryKey,
    initialPageParam: 1,
    queryFn: ({ pageParam }) => fetchCards(deckID, front, sort, pageParam),
    getNextPageParam: (lastPage) => {
      const currentPage = lastPage.metadata.current_page;
      const endPage = lastPage.metadata.last_page;
      if (currentPage >= endPage) {
        return undefined;
      }
      return currentPage + 1;
    },
  });

  const deleteCardMutation = useMutation({
    mutationFn: ({ deckID, cardID }: { deckID: number; cardID: number }) =>
      deleteCard(deckID, cardID),
    onSuccess: (_, variables) => {
      queryClient.setQueryData(queryKey, (oldData: any) => {
        const newPagesArray = oldData.pages.map((page: any) => {
          return {
            ...page,
            cards: page.cards.filter(
              (card: Flashcard) => card.id !== variables.cardID,
            ),
          };
        });
        return {
          pages: newPagesArray,
          pageParams: oldData.pageParams,
        };
      });
      toast.success("Card deleted successfully.");
    },
  });

  // Issue with card changing deck_id to deckID because of Flashcard Type
  const updateCardMutation = useMutation({
    mutationFn: ({
      deckID,
      cardID,
      front,
      back,
    }: {
      deckID: number;
      cardID: number;
      front: string;
      back: string;
    }) => updateCard(deckID, cardID, front, back),
    onSuccess: (_, variables) => {
      queryClient.setQueryData(queryKey, (oldData: any) => {
        const updatedCard: Flashcard = {
          id: variables.cardID,
          deckID: variables.deckID,
          front: variables.front,
          back: variables.back,
        };
        const newPagesArray = oldData.pages.map((page: any) => {
          return {
            ...page,
            cards: page.cards.map((card: any) =>
              card.id === variables.cardID ? updatedCard : card,
            ),
          };
        });
        return {
          pages: newPagesArray,
          pageParams: oldData.pageParams,
        };
      });
      toast.success("Card updated successfully.");
    },
  });

  if (status === "pending") {
    return <p>Loading...</p>;
  }
  if (status === "error") {
    return <p>Error: {error.message}</p>;
  }
  const cardsM = data.pages.flatMap((page) => page.cards) ?? [];
  console.log(data);

  return (
    <div className="grid items-start grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 py-2 px-4">
      {/* {cards && */}
      {/*   cards.map((card) => ( */}
      {/*     <FlashcardItem */}
      {/*       key={card.id} */}
      {/*       card={card} */}
      {/*       deleteCard={deleteCard} */}
      {/*       updateCard={updateCard} */}
      {/*     /> */}
      {/*   ))} */}
      {cardsM &&
        cardsM.map((card) => (
          <FlashcardItem
            key={card.id}
            card={card}
            deleteCard={(cardID) =>
              deleteCardMutation.mutateAsync({ cardID, deckID })
            }
            updateCard={(front: string, back: string) =>
              updateCardMutation.mutateAsync({
                deckID: card.deck_id,
                cardID: card.id,
                front: front,
                back: back,
              })
            }
          />
        ))}
    </div>
  );
}
