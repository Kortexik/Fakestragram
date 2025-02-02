-- +goose Up
CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL, -- The recipient of the notification
    type VARCHAR(50) NOT NULL, -- Type of notification (e.g., 'like', 'comment', 'follow')
    related_user_id INTEGER, -- The user who triggered the notification (optional)
    post_id INTEGER, -- The related post (optional, for likes/comments)
    content TEXT, -- Additional details or a message
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    seen BOOLEAN DEFAULT 0, -- Indicates if the notification has been seen
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (related_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (post_id) REFERENCES user_posts (id) ON DELETE SET NULL
);
-- +goose Down
DROP TABLE notifications;