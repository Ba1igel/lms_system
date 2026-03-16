-- +goose Up
-- +goose StatementBegin
INSERT INTO courses (name, description) VALUES 
('Golang Developer', 'Comprehensive guide to becoming a Go developer. Covers syntax, concurrency, web development, and more.'),
('Python Developer', 'Learn Python from scratch. Perfect for backend developers and data scientists.');

INSERT INTO chapters (name, description, "order", course_id) VALUES
('Introduction to Go', 'Basic setup and understanding of Go', 1, 1),
('Control structures', 'If-else, switch, and loops in Go', 2, 1);

INSERT INTO lessons (name, description, content, "order", chapter_id) VALUES
('Getting Started', 'Installing Go and setting up workspace', 'To install Go, visit golang.org. After installation, set up your GOPATH and run your first Hello World program. Go is a statically typed, compiled programming language designed at Google.', 1, 1),
('If-else Statement in Golang', 'Learning conditional logic', 'In Go, the if statement looks like this: if x > 0 { ... }. You can also declare variables inside the if statement, which makes it very powerful and concise. Example: if err := doSomething(); err != nil { ... }', 1, 2),
('For Loop in Golang', 'The only loop you need', 'Go only has one looping construct, the for loop. The basic for loop has three components separated by semicolons: the init statement, the condition expression, and the post statement.', 2, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM lessons;
DELETE FROM chapters;
DELETE FROM courses;
-- +goose StatementEnd
