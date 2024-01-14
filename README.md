![image](https://logos-world.net/wp-content/uploads/2022/05/MeowChat-Emblem.png)

# Go TCP Chat Server

This is a simple Go TCP chat server that allows multiple clients to connect, chat in various rooms, change usernames, and exit gracefully.

## Features

- Multiple clients can connect simultaneously.
- Clients can join different chat rooms.
- Clients can change their usernames.
- Commands such as /join, /leave, /rn (rename), and /exit are supported.

## Prerequisites

- Go (Golang) installed on your system.

## Usage

1. Clone the repository:

   ```sh
   git clone https://learn.reboot01.com/git/malsamma/net-cat.git
   ```

2. Build the server:

   ```sh
   go build
   ```

3. Start the server:

   ```sh
   ./TCPChat <PORT>
   ```

4. Connect to the server using a TCP client such as `netcat`:

   ```sh
   nc localhost <PORT>
   ```

## Commands

- `/join <room_id>`: Join a chat room (e.g., `/join 1`).
- `/leave`: Leave the current chat room.
- `/rn <new_username>`: Change your username (e.g., `/rn new_username`).
- `/help`: View the list of available commands.
- `/exit`: Disconnect from the server and exit the chat.

## Contributing

If you'd like to contribute to this project, please follow these steps:

1. Fork the project.
2. Create a new branch with your feature or bug fix: `git checkout -b feature/your-feature`.
3. Commit your changes: `git commit -m 'Add some feature'`.
4. Push your changes to your fork: `git push origin feature/your-feature`.
5. Create a pull request on the master repository.

## Contributers
[sahmedG](https://github.com/sahmedG) (Sameer Goumaa)
[MSK17A](https://github.com/MSK17A) (Mohammed AlSammak)
