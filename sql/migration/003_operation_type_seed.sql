-- +goose Up
INSERT INTO operation_types(id, description, mode) VALUES
(1, 'Normal Purchase', 'DEBIT'),
(2, 'Purchase With Installments', 'DEBIT'),
(3, 'Withdrawal', 'DEBIT'),
(4, 'Credit Voucher', 'CREDIT');


-- +goose Down
DELETE FROM operation_types WHERE id IN (1, 2, 3, 4);
