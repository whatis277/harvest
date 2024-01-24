INSERT INTO payment_methods
(id, user_id, label, last4, brand, exp_month, exp_year)
VALUES
-- statics
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0001-000000000001', 'valid', '4242', 'visa', 12, 2028),
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0001-000000000001', 'expired', '5555', 'mastercard', 12, 2022),
-- subscriptions related
('00000000-0000-0000-0001-000000000001', '00000000-0000-0000-0002-000000000001', 'valid', '4242', 'visa', 12, 2028)
;

---- create above / drop below ----

DELETE FROM payment_methods;
