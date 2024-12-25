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
<<<<<<< HEAD
			profile_pic VARCHAR(255) DEFAULT ''
=======
			profilepic VARCHAR(255) DEFAULT ''
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
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
<<<<<<< HEAD
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
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
	`,
=======
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
	// Shops
	`
		CREATE TABLE IF NOT EXISTS shops (
			id SERIAL PRIMARY KEY,
<<<<<<< HEAD
			post_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			location VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			country_id INT NOT NULL,
			avg_rating DECIMAL(2,1) DEFAULT NULL,
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
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
=======
			title VARCHAR(255) NOT NULL,
			description TEXT,
			location VARCHAR(255) NOT NULL,
			country_id INT NOT NULL,
			rating DECIMAL(3,2),
			subforum_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (country_id) REFERENCES countries(id),
			FOREIGN KEY (subforum_id) REFERENCES subforums(id)
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
		);
	`,
	// Ratings
	`
		CREATE TABLE IF NOT EXISTS ratings (
			id SERIAL PRIMARY KEY,
			shop_id INT NOT NULL,
			user_id INT NOT NULL,
<<<<<<< HEAD
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
			url VARCHAR(255) NOT NULL
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
	// Link post and photo
	`
		CREATE TABLE IF NOT EXISTS post_photo (
			post_id INT NOT NULL,
			photo_id INT NOT NULL,
			PRIMARY KEY (post_id, photo_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE
		)
	`,
	// Link comments and photos
	`
		CREATE TABLE IF NOT EXISTS comment_photo (
			comment_id INT NOT NULL,
			photo_id INT NOT NULL,
			PRIMARY KEY (comment_id, photo_id),
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE
		)
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
=======
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
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
}