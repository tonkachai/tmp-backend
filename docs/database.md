## Database ER Diagram

This file contains a Mermaid ER diagram generated from the project's Go models (`models/user.go` and `models/transfer.go`). It shows entities, fields, constraints, and the relationships between them.

```mermaid
erDiagram
    USERS {
        uint id PK "primary key"
        varchar member_code "unique"
        varchar email "unique, not null"
        varchar password "hashed / hidden"
        varchar first_name
        varchar last_name
        varchar phone
        datetime birthday
        datetime created_at
        datetime updated_at
    }

    TRANSFERS {
        uint id PK
        uint from_id FK "references users.id"
        uint to_id FK "references users.id"
        bigint amount
        varchar memo
        datetime created_at
    }

    %% Relationships
    TRANSFERS }o--|| USERS : "from"
    TRANSFERS }o--|| USERS : "to"
```

### Notes
- Entities and fields were derived directly from `models/user.go` and `models/transfer.go`.
- The database is migrated using GORM in `db/db.go` (AutoMigrate for `User` and `Transfer`).
- `from_id` and `to_id` in `TRANSFERS` are foreign keys to `USERS.id` (each transfer has one sender and one receiver; a user can send or receive many transfers).

### Requirements coverage
- Create ER diagram from models: Done (diagram above).
- Output path `docs/database.md`: Done.
