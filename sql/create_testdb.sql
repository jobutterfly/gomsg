CREATE DATABASE gomsg_testing;

CREATE TABLE IF NOT EXISTS boards(
	board_id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL
);

CREATE TABLE threads(
    thread_id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    comment VARCHAR(1275) NOT NULL,
    date VARCHAR(15) NOT NULL,
    board_id INT NOT NULL,
    CONSTRAINT fk_board
    FOREIGN KEY (board_id)
    REFERENCES boards(board_id)
	ON UPDATE CASCADE
	ON DELETE CASCADE
);

CREATE TABLE replies(
    reply_id INT AUTO_INCREMENT PRIMARY KEY,
    comment VARCHAR(1275) NOT NULL,
    date VARCHAR(15) NOT NULL,
    thread_id INT NOT NULL,
    CONSTRAINT fk_thread
    FOREIGN KEY (thread_id)
    REFERENCES threads(thread_id)
	ON UPDATE CASCADE
	ON DELETE CASCADE
);

INSERT INTO boards VALUES (1, "sports"), (2, "random"), (3, "tech");









