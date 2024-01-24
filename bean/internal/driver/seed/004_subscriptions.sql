INSERT INTO subscriptions
(id, user_id, payment_method_id, label, provider, amount, interval, period)
VALUES
-- statics
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0002-000000000001', '00000000-0000-0000-0001-000000000001', 'monthly', 'bean', 500, 1, 'month'),
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0002-000000000001', '00000000-0000-0000-0001-000000000001', 'yearly', 'bean', 5000, 1, 'year')
;

---- create above / drop below ----

DELETE FROM subscriptions;
