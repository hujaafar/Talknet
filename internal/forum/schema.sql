-- Create Users table
CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create Posts table
CREATE TABLE IF NOT EXISTS Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Create Comments table
CREATE TABLE IF NOT EXISTS Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Create Categories table
CREATE TABLE IF NOT EXISTS Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Create Post_Categories table
CREATE TABLE IF NOT EXISTS Post_Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES Categories(id) ON DELETE CASCADE
);

-- Create Likes_Dislikes table
CREATE TABLE IF NOT EXISTS Likes_Dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    like_dislike BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comments(id) ON DELETE CASCADE
);

-- Create Sessions table
CREATE TABLE IF NOT EXISTS Sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO Categories (name) VALUES ('Technology');
INSERT OR IGNORE INTO Categories (name) VALUES ('Sport');
INSERT OR IGNORE INTO Categories (name) VALUES ('Science');
INSERT OR IGNORE INTO Categories (name) VALUES ('Education');
INSERT OR IGNORE INTO Categories (name) VALUES ('Gaming');
INSERT OR IGNORE INTO Categories (name) VALUES ('TV');
INSERT OR IGNORE INTO Categories (name) VALUES ('Comedy');
INSERT OR IGNORE INTO Categories (name) VALUES ('History');
INSERT OR IGNORE INTO Categories (name) VALUES ('Social');
INSERT OR IGNORE INTO Categories (name) VALUES ('Finance');
INSERT OR IGNORE INTO Categories (name) VALUES ('News');
INSERT OR IGNORE INTO Categories (name) VALUES ('Others');

-- Support the feed, profile, reply and session lookups without scanning every row.
-- Additive tables keep existing community databases compatible.
CREATE TABLE IF NOT EXISTS Bookmarks (
    user_id INTEGER NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    post_id INTEGER NOT NULL REFERENCES Posts(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id)
);
CREATE TABLE IF NOT EXISTS Post_Revisions (
    post_id INTEGER PRIMARY KEY REFERENCES Posts(id) ON DELETE CASCADE,
    revision INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_posts_user ON Posts(user_id, id);
CREATE INDEX IF NOT EXISTS idx_comments_post ON Comments(post_id, id);
CREATE INDEX IF NOT EXISTS idx_post_categories_post ON Post_Categories(post_id, category_id);
CREATE INDEX IF NOT EXISTS idx_post_categories_category ON Post_Categories(category_id, post_id);
CREATE INDEX IF NOT EXISTS idx_reactions_post ON Likes_Dislikes(post_id, user_id, like_dislike);
CREATE INDEX IF NOT EXISTS idx_reactions_comment ON Likes_Dislikes(comment_id, user_id, like_dislike);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON Sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON Sessions(expires_at);

-- Records explicit demo-pack imports so a repeated command cannot duplicate them.
CREATE TABLE IF NOT EXISTS Demo_Seeds (
    name TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
