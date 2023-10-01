
# Project Name

This project provides an easy way to label bank transactions, making it easier to manage
personal finances. It leverages on Nordigen api through ricardograndecros/go-nordigen 
client. 

The web application consists of a simple frontend where transactions are listed. Each
transaction contains a dropdown list with predefined labels by the user. 

## ROADMAP
1. Redesign frontend for a more appealing look
2. Add refresh button to the frontend
3. Add send button to transaction cards in the frontend
4. Add send all labeled transactions button to the frontend
5. Generalize data export via config files so that it is possible to send labeled transactions to any desired endpoint


## Tech Stack

- **Backend:** Go
  - **Database:** SQLite (temporarily)
- **Frontend:** Preact

## Prerequisites

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/en/download/) & [npm](https://www.npmjs.com/get-npm) (for the Preact frontend)

## Setup & Installation

### Backend

1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```

2. Install the required Go packages:
   ```bash
   go get -v ./...
   ```

3. Run the Go server:
   ```bash
   go run server.go
   ```

### Frontend

1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```

2. Install the required npm packages:
   ```bash
   npm install
   ```

3. Start the Preact development server:
   ```bash
   npm run start
   ```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

This is just a project I am building for learning and for fun. I understand it needs lots of refactoring. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
