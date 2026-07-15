CREATE TABLE qr_configs (
    id SERIAL PRIMARY KEY,
    link_id INT REFERENCES links(id),
    foreground_color VARCHAR(7) DEFAULT '#000000',
    background_color VARCHAR(7) DEFAULT '#FFFFFF',
    dot_style VARCHAR(20) DEFAULT 'square',
    eye_style VARCHAR(20) DEFAULT 'square',
    logo_url TEXT,
    size INT DEFAULT 1024,
    updated_at TIMESTAMP DEFAULT NOW()
);