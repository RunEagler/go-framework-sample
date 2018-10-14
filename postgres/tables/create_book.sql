CREATE TABLE book (
    book_id SERIAL PRIMARY KEY,
    title text NOT NULL,
    pages numeric(5,0) NOT NULL,
    date_of_publication DATE
);
