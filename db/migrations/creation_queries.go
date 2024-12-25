package migrations

var createTablesQueries = [...]string{
	// Users
	`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			profilepic VARCHAR(255) DEFAULT ''
		);
	`,
	// Countries
	`
		CREATE TABLE IF NOT EXISTS countries(
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT
		);
	`,
	// Subforums
	`
		CREATE TABLE IF NOT EXISTS subforums (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT
		);
	`,
	// Shops
	`
		CREATE TABLE IF NOT EXISTS shops (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			location VARCHAR(255) NOT NULL,
			country_id INT NOT NULL,
			rating DECIMAL(3,2),
			subforum_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (country_id) REFERENCES countries(id),
			FOREIGN KEY (subforum_id) REFERENCES subforums(id)
		);
	`,
	// Ratings
	`
		CREATE TABLE IF NOT EXISTS ratings (
			id SERIAL PRIMARY KEY,
			shop_id INT NOT NULL,
			user_id INT NOT NULL,
			rating DECIMAL(2,1) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (shop_id) REFERENCES shops(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`,
	// Posts
	`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY, 
			subforum_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (subforum_id) REFERENCES subforums(id)
		);
	`,
	// Comments
	`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL,
			user_id INT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`,
	// Post_Country join table
	`
		CREATE TABLE IF NOT EXISTS post_country (
			post_id INT NOT NULL,
			country_id INT NOT NULL,
			PRIMARY key (post_id, country_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (country_id) REFERENCES countries(id)
		);
	`,
}