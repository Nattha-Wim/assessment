

-- Table Definition
CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);

INSERT INTO expenses (id, title, amount, note, tags) VALUES (1, 'apple juice', 75.00, 'cafe', ARRAY['drink']);