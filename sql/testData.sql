-- admin (password - '12345678')
INSERT INTO users (username, password, is_admin) VALUES 
    ('mr. Fedya', '$2a$10$W1uTjnpz.h/hbfWuRhO04ekfs6FffeMsIbtFpxLiFhE6eMgW7oMUi', TRUE)
ON CONFLICT (username) DO NOTHING
RETURNING id, username, password, is_admin, active, created;

-- default admins token
INSERT INTO users_tokens (token, user_id) VALUES 
    ('defaultAdminsToken', 1);
