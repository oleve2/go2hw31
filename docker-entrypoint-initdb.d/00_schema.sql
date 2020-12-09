-- табличка с пользователями и их паролями
CREATE TABLE users
(
    id       BIGSERIAL PRIMARY KEY,
    login    TEXT      NOT NULL UNIQUE,
    password TEXT      NOT NULL,
    roles    TEXT[]    NOT NULL DEFAULT '{}',
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- табличка с токенами (чтобы пользователь каждый раз не присылал пароль и логин)
CREATE TABLE tokens (
    id TEXT PRIMARY KEY,
    userId BIGINT NOT NULL REFERENCES users,
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- таблца карт
CREATE TABLE cards(
    id              BIGSERIAL PRIMARY KEY, -- записывайте UUID - автоинкремент
    type            TEXT NOT NULL,
    bank_name       TEXT NOT NULL,
    card_number     TEXT NOT NULL,
    card_due_date   TEXT NOT NULL,
    balance         INTEGER NOT NULL,
    user_id         INTEGER NOT NULL REFERENCES users,
    is_virtual      boolean NOT NULL
);

