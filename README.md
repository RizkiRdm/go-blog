# GoBlog app

Welcome to our GoBlog app! This is a simple web application where users can create, read, update, and delete blog posts.

## Features

- **User Authentication**: Users can sign up, log in, and log out securely.
- **Create Blog Posts**: Authenticated users can create new blog posts with a title, content, and an optional image thumbnail.
- **Read Blog Posts**: Users can view all blog posts or search for specific posts using keywords.
- **Update Blog Posts**: Users can edit their own blog posts, including changing the title, content, or thumbnail image.
- **Delete Blog Posts**: Users can delete their own blog posts if they no longer wish to keep them.
- **Search Blog**: Users can search blog based on title or body blog.
## Getting Started

To run this application locally, follow these steps:

1. Clone this repository to your local machine.
2. Install dependencies by running `go mod tidy`.
3. Set up the database according to the provided SQL scripts.
4. Configure the environment variables for database connection.
5. Run the application using `go run main.go`.
6. Access the application in your web browser at `http://localhost:3000`.

## Dependencies

- Fiber: A web framework for Go, providing fast and flexible routing.
- UUID: A package for generating universally unique identifiers.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
