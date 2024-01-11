package server

import (
	"log"
	"os"
)

/* function will take ChatMsg struct as an input and displays a message to clients */
func Broadcast(message ChatMsg) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	logFile, err := os.OpenFile("chat_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	_, err = logFile.WriteString(message.Content + "\n")
	if err != nil {
		log.Println("Error writing to log file:", err)
		return
	}

	chatRooms[message.RoomID] = append(chatRooms[message.RoomID], message)

	for client := range clients {
		if client.RoomID == message.RoomID && client.no_get == false {
			client.Writer.WriteString(formatMessage(message) + "\n")
			client.Writer.Flush()
		}
	}
}

func Renaming_Arg_check(args []string) bool {
	return len(args) == 2
}
