INSERT INTO "users" ("id", "email", "password", "created_at", "updated_at") VALUES(123, '123@test.com', '$2a$10$9J3sIOgWlMVvssEEmoUm.eBHKembea4CLBqwHfjln4vHfbKOOSdJK', '2024-01-31 15:26:31.804593+00', '2024-01-31 15:26:31.804593+00');

INSERT INTO "articles" ("id", "user_id", "title", "content", "created_at", "updated_at") VALUES (999, 123, 'seeded title 999', 'seeded content 999', '2024-01-31 15:26:31.804593+00', '2024-01-31 15:26:31.804593+00');
INSERT INTO "articles" ("id", "user_id", "title", "content", "created_at", "updated_at") VALUES (888, 123, 'seeded title 888', 'seeded content 888', '2024-01-31 15:26:31.804593+00', '2024-01-31 15:26:31.804593+00');
