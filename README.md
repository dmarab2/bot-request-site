# bot-request-site
This is a request site for chatbot creators and users to view and create requests for chatbots.
This application is made up of two distinct parts: a frontend written in React using Typescript, and a backend written written in Go.
The backend is written mainly using the standard net/http library and is connected to a Postgres database via usage of sqlc. The goose package is also used for schema building.
The frontend is written in Next.js using Typescript.


TODO:
Do some refactoring in main to break the logic out into their own functions
Related to the above, write tests for logic functions
Edit autocomplete query to use postgres trigram extension
check tag aliases schema I feel like I'm forgetting something
Dockerize the backend and frontend with a PostgreSQL image as well.