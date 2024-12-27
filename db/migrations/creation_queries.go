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
			profile_pic VARCHAR(255) DEFAULT ''
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
	// Posts
	`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY, 
			subforum_id INT NOT NULL,
			user_id INT,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			comment_count INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (subforum_id) REFERENCES subforums(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`,
	// Shops
	`
		CREATE TABLE IF NOT EXISTS shops (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			avg_rating DECIMAL(2,1) DEFAULT NULL,
			country_id INT NOT NULL,
			type VARCHAR(255) NOT NULL,
			lat DOUBLE PRECISION NOT NULL,
			long DOUBLE PRECISION NOT NULL,
			address VARCHAR(255),
			map_link VARCHAR(255) NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
		);
	`,
	// Comments
	`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL,
			user_id INT,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`,
	// Ratings
	`
		CREATE TABLE IF NOT EXISTS ratings (
			id SERIAL PRIMARY KEY,
			shop_id INT NOT NULL,
			user_id INT NOT NULL,
			score DECIMAL(1,0) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT unique_user_shop_rating UNIQUE (shop_id, user_id) -- To ensure one user one rating per shop --
		);
	`,
	// Photos
	`
		CREATE TABLE IF NOT EXISTS photos (
			id SERIAL PRIMARY KEY,
			url VARCHAR(255) NOT NULL,
			context VARCHAR(255) NOT NULL,
			context_id INT NOT NULL
		)
	`,
	// Link post and country
	`
		CREATE TABLE IF NOT EXISTS post_country (
			post_id INT NOT NULL,
			country_id INT NULL,
			PRIMARY key (post_id, country_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
		);
	`,
	// Trigger to update average rating in 'shop_posts' table
	`
		CREATE OR REPLACE FUNCTION update_average_rating()
		RETURNS TRIGGER AS $$
		BEGIN
			UPDATE shops
			SET avg_rating = (
				SELECT AVG(score)
				FROM ratings
				WHERE shop_id = COALESCE(NEW.shop_id, OLD.shop_id)
			)
			WHERE id = COALESCE(NEW.shop_id, OLD.shop_id);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql
	`,
	`
		CREATE OR REPLACE TRIGGER after_rating_insert
		AFTER INSERT ON ratings
		FOR EACH ROW
		EXECUTE FUNCTION update_average_rating();
	`,
	`
		CREATE OR REPLACE TRIGGER after_rating_update
		AFTER UPDATE ON ratings
		FOR EACH ROW
		EXECUTE FUNCTION update_average_rating();
	`,
	`
		CREATE OR REPLACE TRIGGER after_rating_delete
		AFTER DELETE ON ratings
		FOR EACH ROW
		EXECUTE FUNCTION update_average_rating();
	`,
	// trigger to update comment count
	`
		CREATE OR REPLACE FUNCTION update_comment_count()
		RETURNS TRIGGER AS $$
		BEGIN 
			UPDATE posts
			SET comment_count = (
				SELECT COUNT(*)
				FROM comments
				WHERE post_id = NEW.post_id
			)
			WHERE id = NEW.post_id;

			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`,
	`
		CREATE OR REPLACE TRIGGER after_comment_insert
		AFTER INSERT ON comments
		FOR EACH ROW
		EXECUTE FUNCTION update_comment_count();
	`,
	`
		CREATE OR REPLACE TRIGGER after_comment_delete
		AFTER DELETE ON comments
		FOR EACH ROW
		EXECUTE FUNCTION update_comment_count();
	`,
	// indexes
	`
		CREATE INDEX IF NOT EXISTS idx_posts_title
		ON posts(title);
	`,
	`
		CREATE INDEX IF NOT EXISTS idx_comments_post_id
		ON comments(post_id)
	`,
	`
		CREATE INDEX IF NOT EXISTS idx_ratings_shop_post_id
		ON ratings(shop_id)
	`,
}
