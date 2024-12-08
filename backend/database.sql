CREATE TABLE Account (
	acc_id SERIAL PRIMARY KEY,
	username VARCHAR(50) NOT NULL,
	email VARCHAR(100) NOT NULL
  );
  
CREATE TABLE Character (
	char_id SERIAL PRIMARY KEY,
	acc_id INT REFERENCES Account(acc_id) ON DELETE CASCADE,
	class_id INT NOT NULL
  );
  
CREATE TABLE Scores (
	score_id SERIAL PRIMARY KEY,
	char_id INT REFERENCES Character(char_id) ON DELETE CASCADE,
	reward_score INT NOT NULL
  );
  