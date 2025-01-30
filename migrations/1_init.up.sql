CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       access_level INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_username ON users(username);

CREATE TABLE tasks (
                       task_id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       status VARCHAR(50) DEFAULT 'todo',
                       deadline TIMESTAMP,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE task_assignments (
                                  id SERIAL PRIMARY KEY,
                                  user_id INT NOT NULL REFERENCES users(user_id),
                                  task_id INT NOT NULL REFERENCES tasks(task_id),
                                  UNIQUE(user_id, task_id)
);