CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_type VARCHAR(20) NOT NULL,
    amount INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    receipt_type VARCHAR(20),
    receipt_file_id TEXT,
    receipt_text TEXT,
    reviewed_by_tg_id BIGINT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_admin_msgs (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    admin_tg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    message_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(payment_id, admin_tg_id)
);
