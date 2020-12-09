package business

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// -----------------------------------------------------------------------
/* go2hw31 */

type Card struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	BankName   string `json:"bank_name"`
	CardNumber string `json:"card_number"`
	Balance    int64  `json:"balance"`
	UserID     int64  `json:"user_id"`
}

func (s *Service) GetAllCards(ctx context.Context) ([]*Card, error) {
	cards := make([]*Card, 0)

	rows, err := s.pool.Query(ctx, `select id, type, bank_name, card_number, balance, user_id from cards order by id asc`)
	if err != nil {
		if err == pgx.ErrNoRows {
			return cards, nil
		}
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		cardItem := &Card{}
		err = rows.Scan(
			&cardItem.ID,
			&cardItem.Type,
			&cardItem.BankName,
			&cardItem.CardNumber,
			&cardItem.Balance,
			&cardItem.UserID,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		cards = append(cards, cardItem)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cards, nil
}

func (s *Service) GetUserCards(ctx context.Context, id int64) ([]*Card, error) {
	cards := make([]*Card, 0)

	rows, err := s.pool.Query(
		ctx,
		`select id, type, bank_name, card_number, balance, user_id from cards where user_id=$1 order by id asc`,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return cards, nil
		}
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		cardItem := &Card{}
		err = rows.Scan(
			&cardItem.ID,
			&cardItem.Type,
			&cardItem.BankName,
			&cardItem.CardNumber,
			&cardItem.Balance,
			&cardItem.UserID,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		cards = append(cards, cardItem)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cards, nil
}

// -----------------------------------------------------------------------
