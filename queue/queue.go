package queue

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"pratyushtiwary/sqs/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type QueueConfig struct {
	BufferSize int
	Timeout    int
}

type Queue struct {
	DB     *gorm.DB
	Config QueueConfig
}

func Init(config QueueConfig) (*Queue, error) {
	db, err := gorm.Open(sqlite.Open("queue.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	models.Init(db)

	queue := new(Queue)

	queue.DB = db
	queue.Config = config

	return queue, nil
}

func (q *Queue) HandleRequest(conn net.Conn) {
	defer conn.Close()

	q.readData(conn)
}

func (q *Queue) readData(conn net.Conn) {
	buffer := make([]byte, q.Config.BufferSize)

	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Duration(q.Config.Timeout) * time.Second))
		if err != nil {
			log.Println("SetReadDeadline failed:", err)
			return
		}

		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF && !errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("Error:", err)
			}
			break
		}
		if n == 0 {
			break
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}
