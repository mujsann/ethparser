# Documentation

## Overview

This application provides a simple interface to interact with the Ethereum blockchain; will allow users to retrieve block information, subscribe to address notifications, and fetch transaction details.

The `parser` package (/pkg) provides functionality to interact with the Ethereum blockchain. It provides the methods for fetching block data, processing transactions, and validating Ethereum addresses.

*Note: This application does not use any external Go libraries.

## Requirements

- Go 1.16 or higher

The application requires the following environment variables to be set:

- `TRANSACTION_DAYS_LIMIT`: The number of days to limit the transactions fetched.
- `TEST_ADDRESS`: A valid Ethereum address for testing purposes.
- `PORT`: The port on which the server will run (default is 8080).

### Example `.env` file

```text
TRANSACTION_DAYS_LIMIT=5
TEST_ADDRESS=0x20485641350a3ca182d84199b7b3f679f03703bf
PORT=8080
```

## Installation

1. Clone the repository:

   ```bash
   git clone git@github.com:mujsann/ethparser.git
   cd <repository-directory>
   ```

2. Create a `.env` file in the root directory and set the required environment variables.

3. Build the application:

   ```bash
   go build
   ```

4. Run the application:

   ```bash
   ./ethparser
   ```

5. Build and run the application at once

    ```bash
    ./run_app
    ```

## API Endpoints

| Method | Endpoint                | Description                          |
|--------|-------------------------|--------------------------------------|
| GET    | /current-block          | Retrieves the current block number.  |
| POST   | /subscribe              | Subscribes an address for notifications. |
| GET    | /transactions/{address} | Fetches transactions for a specified address. |

### 1. Subscribe to Address Updates

- **Endpoint**: `/subscribe`
- **Method**: `POST`
- **Request Body**:

  ```json
  {
      "address": "<Ethereum Address>"
  }
  ```

- **Response**:
  - `200 OK` if the subscription is successful.
  - `400 Bad Request` if the address is invalid or already subscribed.

### 2. Get Current Block

- **Endpoint**: `/current-block`
- **Method**: `GET`
- **Response**:
  - `200 OK` with the current block number in the response body:

  ```json
  {
      "current_block": <block_number>
  }
  ```

### 3. Get Transactions for an Address

- **Endpoint**: `/transactions/{address}`
- **Method**: `GET`
- **Response**:
  - `200 OK` with a list of transactions for the specified address.
  - `400 Bad Request` if the address is invalid.

## Testing

- **Run Tests**: Execute the following command to run the tests:

  ```bash
  go test ./...
  ```
