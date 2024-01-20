INSERT INTO login_tokens
(id, email, hashed_token, created_at, expires_at)
VALUES
-- statics
('00000000-0000-0000-0000-000000000001', 'static-expired', 'hashed-token', '2024-01-21 00:00:00', '2024-01-21 00:10:00'),
-- actionables
('10000000-0000-0000-0000-000000000000', 'action-overwrite', 'old-hashed-token', '2024-01-21 00:00:00', '2024-01-21 00:10:00'),
('20000000-0000-0000-0000-000000000000', 'action-delete', 'hashed-token', '2024-01-21 00:00:00', '2024-01-21 00:10:00')
;