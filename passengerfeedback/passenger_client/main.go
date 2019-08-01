package main

import (
	pb "../proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
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
				fmt.Println("could not add passenger:", err)
			} else {
				fmt.Println(r.Msg)
				if r.ErrorCode == pb.Error_SUCCESS {
					fmt.Println(r.Data.String())
				}
			}

			wg.Done()
		}
	}()

	wg.Wait()

	fmt.Println("------------------------------------------------")
	fmt.Println("----------GET FEEDBACK BY PASSENGER ID----------")

	feedBackByPassenger, err := c.GetFeedbackByPassengerId(ctx, &pb.GetPassengerFeedbackByPassengerIdRequest{
		PassengerId: 1,
		Offset: 0,
		Limit: 2,
	})
	if err != nil {
		fmt.Println("could not get feedbacks by passenger:", err)
	} else {
		fmt.Println(feedBackByPassenger.Msg)

		if feedBackByPassenger.ErrorCode == pb.Error_SUCCESS && len(feedBackByPassenger.Data) > 0 {
			for _, v := range feedBackByPassenger.Data {
				fmt.Println(v)
			}
		}
	}

}
