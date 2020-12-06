
INSERT INTO users(login, password, roles) VALUES
('admin', '$2a$10$ctPFhgJh.YIE21AA0OGl5er3p9f3XsAwkmTXnP2I7BxCpQbr1QAg2', '{"ADMIN", "USER"}'), -- у этого пользователя две роли (т.е. он и админ, и обычный юзер)
('user', '$2a$10$ctPFhgJh.YIE21AA0OGl5er3p9f3XsAwkmTXnP2I7BxCpQbr1QAg2', '{"USER"}');

insert into payments (senderId, amount, comment) VALUES
(1,100,'comment 1-100'), 
(1,200,'comment 1-200'), 
(1,300,'comment 1-300'), 
(2,100,'comment 2-100'), 
(2,200,'comment 2-200');

insert into cards (type, bank_name, card_number, card_due_date, balance, user_id, is_virtual) VALUES
('visa','Tinkoff','1111-2222-3333-0001', '2023-01-04', 1000, 1, FALSE),
('visa','Tinkoff','1111-2222-3333-0002', '2023-01-04', 2000, 1, FALSE),

('visa','Tinkoff','1111-2222-3333-0003', '2023-01-04', 3000, 2, FALSE),
('visa','Tinkoff','1111-2222-3333-0004', '2023-01-04', 4000, 2, FALSE),
('visa','Tinkoff','1111-2222-3333-0005', '2023-01-04', 5000, 3, FALSE);


