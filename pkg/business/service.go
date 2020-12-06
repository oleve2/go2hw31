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

/*

// -----------------------------------------------------------------------
type Payment struct {
	ID       int64  `json:"id"`
	SenderID int64  `json:"sender_id"`
	Amount   int64  `json:"amount"`
	Comment  string `json:"comment"`
}

func (s *Service) GetAllPayments(ctx context.Context) ([]*Payment, error) {
	payments := make([]*Payment, 0)

	rows, err := s.pool.Query(ctx, `select id, senderid, amount, comment from payments order by id asc`)
	if err != nil {
		if err == pgx.ErrNoRows {
			return payments, nil
		}
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		paymentItem := &Payment{}
		err = rows.Scan(&paymentItem.ID, &paymentItem.SenderID, &paymentItem.Amount, &paymentItem.Comment)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		payments = append(payments, paymentItem)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return payments, nil
}

func (s *Service) GetUserPayments(ctx context.Context, id int64) ([]*Payment, error) {
	payments := make([]*Payment, 0)

	rows, err := s.pool.Query(
		ctx,
		`select id, senderid, amount, comment from payments where senderid=$1 order by id asc`,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return payments, nil
		}
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		paymentItem := &Payment{}
		err = rows.Scan(&paymentItem.ID, &paymentItem.SenderID, &paymentItem.Amount, &paymentItem.Comment)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		payments = append(payments, paymentItem)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return payments, nil
}

func (s *Service) PostUserPayments(ctx context.Context, senderID int64, amount int64, comment string) error {
	tag, err := s.pool.Exec(
		ctx,
		`insert into payments (senderid, amount, comment) values ($1, $2, $3)`,
		senderID, amount, comment,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("No rows updated")
	}

	return nil
}

*/
