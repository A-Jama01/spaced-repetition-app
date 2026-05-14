import {
  Field,
  FieldGroup,
  FieldSet,
  FieldLabel,
  FieldDescription,
  FieldTitle,
} from "./ui/field";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useState } from "react";
import type { Flashcard } from "@/pages/Home";
import { ScrollArea } from "./ui/scroll-area";

interface CardFormProps {
  title: string;
  operationName: string;
  card?: Flashcard;
  children: React.ReactNode;
  cardOperation(front: string, back: string): Promise<void>;
}

export default function CardForm({
  title,
  operationName,
  card,
  children,
  cardOperation,
}: CardFormProps) {
  const initialFront = card === undefined ? "" : card.front;
  const initalBack = card === undefined ? "" : card.back;
  const [frontInput, setFrontInput] = useState<string>(initialFront);
  const [backInput, setBackInput] = useState<string>(initalBack);
  const [showPreview, setShowPreview] = useState<boolean>(false);

  async function handleSubmit() {
    // TODO: create may be broken after refactor. erase when complete.
    try {
      await cardOperation(frontInput, backInput);
      if (operationName === "Create") {
        setFrontInput("");
        setBackInput("");
      }
      setShowPreview(false);
    } catch (err) {
      return;
    }
  }

  return (
    <form action={handleSubmit}>
      <FieldSet>
        <FieldTitle>{title}</FieldTitle>
        <FieldDescription>
          Fill in front and back. Use Markdown for formatting.
        </FieldDescription>
        <FieldGroup className="overflow-auto">
          <Field>
            <FieldLabel htmlFor="front">Front</FieldLabel>
            {showPreview ? (
              <div className="prose prose-sm border-2 rounded-md p-2 shadow-sm">
                <Markdown remarkPlugins={[remarkGfm]}>{frontInput}</Markdown>
              </div>
            ) : (
              <Textarea
                rows={12}
                maxLength={300}
                value={frontInput}
                onChange={(e) => setFrontInput(e.currentTarget.value)}
              />
            )}
          </Field>
          <Field>
            <FieldLabel htmlFor="back">Back</FieldLabel>
            {showPreview ? (
              <div className="prose prose-sm border-2 rounded-md p-2 shadow-sm">
                <Markdown remarkPlugins={[remarkGfm]}>{backInput}</Markdown>
              </div>
            ) : (
              <Textarea
                rows={12}
                maxLength={2000}
                value={backInput}
                onChange={(e) => setBackInput(e.currentTarget.value)}
              />
            )}
          </Field>
        </FieldGroup>
        <Field>
          <div className="flex justify-between">
            <Button
              type="button"
              variant={"outline"}
              onClick={() => setShowPreview(!showPreview)}
            >
              Preview
            </Button>
            <div className="grid grid-cols-2 gap-2">
              <Button type="submit">{operationName}</Button>
              {children}
            </div>
          </div>
        </Field>
      </FieldSet>
    </form>
  );
}
