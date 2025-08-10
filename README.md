# Book Library Assessment

- **Server docs:** [`server/README.md`](./server/README.md)
- **Client docs:** [`client/README.md`](./client/README.md)

## Quick Start

### 1) Server

```bash
cd server
# Optional: spin Postgres in Docker and run API with migrations
make run/dev
# or run just the API if DB is already up
make run
```

### 2) Client

```bash
cd client
# Create `client/.env.local` with:
#API_BASE_URL=http://localhost:4748
# (ensure API_BASE_URL=http://localhost:4748)
npm install
npm run dev
```

## Screenshots

> The images below live in [`/screenshots`](./screenshots). Some have “2” variants when multiple captures were helpful.

### Backend (Swagger)

![Swagger Overview](./screenshots/swagger.png)
![Create Book](./screenshots/swagger-create-book.png)
![Create Book (alt)](./screenshots/swagger-create-book2.png)
![Get Book](./screenshots/swagger-get-book.png)
![URL Process](./screenshots/swagger-url-process.png)
![URL Process (alt)](./screenshots/swagger-url-process2.png)
![Validation 422](./screenshots/swagger-validation-422.png)
![Validation 422 (alt)](./screenshots/swagger-validation-4222.png)

### Frontend

![Books List](./screenshots/books-list.png)
![Add Book (validation)](./screenshots/books-add-validation.png)
![Book Detail](./screenshots/book-detail.png)
![Edit Book](./screenshots/books-edit.png)
![Delete Book](./screenshots/books-delete.png)
![404 Not Found](./screenshots/books-404.png)
![Network Error](./screenshots/books-network-error.png)
