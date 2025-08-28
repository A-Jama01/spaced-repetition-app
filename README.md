# Spaced Repetition Flashcards

A flashcard web application that helps learners retain information more effectively using 
**spaced repetition**. The app implements the latest version of the **Free Spaced Repetition 
Scheduler (FSRS)**, which reduces review load by **20–30% compared to traditional algorithms**.  

Learn more about FSRS here: [ABCs of FSRS](https://github.com/open-spaced-repetition/fsrs4anki/wiki/ABC-of-FSRS).

---

## Features

### Implemented
- Review due flashcards
- FSRS-based scheduling of reviews
- Create, update, and delete decks and cards
- Daily stats with heatmap and due cards forecast
- JWT-based authentication

### Planned
- Card & Deck search/filtering
- React & Tailwind Frontend
- Convert notes into flashcards using AI  
- Public repository for users to share decks  

---

## Tech Stack
- **Backend:** Go  
- **Database:** PostgreSQL  
- **Authentication:** JWT  
- **Frontend:** (planned) React & Tailwind 

## API Overview

Base URL: `/v1`

### Auth
| Method | Endpoint          | Description              |
|--------|------------------|--------------------------|
| POST   | `/auth/register` | Register a new user      |
| POST   | `/auth/login`    | Authenticate and get JWT |

---

### Decks
| Method | Endpoint              | Description            |
|--------|-----------------------|------------------------|
| GET    | `/decks`              | List all decks         |
| POST   | `/decks`              | Create a new deck      |
| PUT    | `/decks/{deck_id}`    | Update a deck          |
| DELETE | `/decks/{deck_id}`    | Delete a deck          |

---

### Cards
| Method | Endpoint                                 | Description                     |
|--------|------------------------------------------|---------------------------------|
| GET    | `/decks/{deck_id}/cards`                 | List all cards in a deck        |
| GET    | `/decks/{deck_id}/cards/due`             | List due cards for review       |
| POST   | `/decks/{deck_id}/cards`                 | Create a new card in a deck     |
| PATCH  | `/decks/{deck_id}/cards/{card_id}`       | Update a card                   |
| DELETE | `/decks/{deck_id}/cards/{card_id}`       | Delete a card                   |
| PATCH  | `/decks/{deck_id}/cards/{card_id}/review`| Review a card (FSRS scheduling) |

---

###  Stats
| Method | Endpoint   | Description                  |
|--------|------------|------------------------------|
| GET    | `/stats`   | Get daily stats, Heatmap, & forecasts  |
