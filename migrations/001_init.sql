CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'client',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    verification_status VARCHAR(20) NOT NULL DEFAULT 'none',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    business_name VARCHAR(255),
    city VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS executor_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(255),
    category VARCHAR(20),
    experience_years INT DEFAULT 0,
    description TEXT,
    portfolio_links TEXT,
    rating FLOAT DEFAULT 0,
    total_orders INT DEFAULT 0,
    completed_orders INT DEFAULT 0,
    response_speed FLOAT DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_pro BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    video_file_id VARCHAR(500),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    reviewed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(300),
    category VARCHAR(20),
    budget_type VARCHAR(20),
    budget_from BIGINT,
    budget_to BIGINT,
    deadline TIMESTAMP,
    refs TEXT,
    is_urgent BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'open',
    max_responses INT DEFAULT 15,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS responses (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    executor_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    proposed_price BIGINT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(task_id, executor_id)
);

CREATE TABLE IF NOT EXISTS task_assignments (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT UNIQUE REFERENCES tasks(id) ON DELETE CASCADE,
    executor_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT UNIQUE REFERENCES tasks(id) ON DELETE CASCADE,
    client_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    executor_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20),
    start_date TIMESTAMP DEFAULT NOW(),
    end_date TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS usage_limits (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    free_responses_left INT DEFAULT 5,
    last_reset_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50),
    payload TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS action_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    action VARCHAR(255),
    data TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
