package main

import (
	"context"
	"fmt"
	pb "github.com/vanlong2712/go-grab-training/passengerfeedback/proto"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	address = "localhost:50052"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
	}
	defer conn.Close()
	c := pb.NewFeedbackClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	number := 10
	wg := sync.WaitGroup{}
	wg.Add(number)

	go func() {
		for i := 1; i <= number; i++ {
			r, err := c.AddPassengerFeedback(ctx, &pb.PassengerFeedbackRequest{
				BookingCode: "DRIVE0" + strconv.Itoa(i),
				PassengerId: 1,
				Feedback:    "Good service.",
			})

			if err != nil {
				log.Fatalf("could not add passenger: %v", err)
			} else {
				fmt.Println(r.Msg)
				if r.ErrorCode == pb.Error_SUCCESS {
					fmt.Printf("The booking code: %v, the passenger id: %d, the feedback: %v\n", r.Data.BookingCode, r.Data.PassengerId, r.Data.Feedback)
				}
			}

			wg.Done()
		}
	}()

	wg.Wait()

}
