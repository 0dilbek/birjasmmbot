CREATE TABLE IF NOT EXISTS verification_admin_msgs (
    id BIGSERIAL PRIMARY KEY,
    verification_id BIGINT NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    admin_tg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    message_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(verification_id, admin_tg_id)
);
