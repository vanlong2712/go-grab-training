package main

import (
	pb "../proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	address = "localhost:50052"
)

var methods = map[int]string {
	1: "Add Feedback",
	2: "Get Feedback By PassengerID",
	3: "Get Feedback By BookingCode",
	4: "Delete Feedback",
	5: "Stop the program",
}

func addAndPrintPassengerFeedback(c pb.FeedbackClient, bookingCode string, passengerId int32, feedback string, printMsg bool) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.AddPassengerFeedback(ctx, &pb.PassengerFeedbackRequest{
		BookingCode: bookingCode,
		PassengerId: passengerId,
		Feedback:    feedback,
	})

	if err != nil {
		fmt.Println("could not add passenger:", err)
		return
	}

	if printMsg {
		fmt.Println(r.Msg)
	}

	if r.ErrorCode == pb.Error_SUCCESS {
		fmt.Println(r.Data.String())
	}

}

func generateSomeFeedback(c pb.FeedbackClient , out chan<- bool, number int) {
	for i := 1; i <= number; i++ {
		addAndPrintPassengerFeedback(c, "DRIVE0" + strconv.Itoa(i), 1, "Good service.", false)
	}
	out <- true
}

func addPassengerFromInput(c pb.FeedbackClient) {
	var passengerId int32
	var bookingCode string
	var feedback string

	fmt.Print("Input Passenger Id: ")
	_, err := fmt.Scan(&passengerId)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Print("Input Booking Code: ")
	_, err = fmt.Scan(&bookingCode)
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(strings.TrimSpace(bookingCode)) == 0 {
		fmt.Println("Booking code expected at least one character.")
		return
	}
	fmt.Print("Input feedback: ")
	_, err = fmt.Scan(&feedback)
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(strings.TrimSpace(feedback)) == 0 {
		fmt.Println("Feedback expected at least one character.")
		return
	}
	addAndPrintPassengerFeedback(c, bookingCode, passengerId, feedback, true)
}

func getFeedbackByPassengerIdFromInput(c pb.FeedbackClient) {
	var passengerId int32
	var offset int32
	var limit int32
	fmt.Print("Input Passenger Id: ")
	_, err := fmt.Scan(&passengerId)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Print("Input offset: ")
	_, err = fmt.Scan(&offset)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Print("Input limit: ")
	_, err = fmt.Scan(&limit)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetFeedbackByPassengerId(ctx, &pb.GetPassengerFeedbackByPassengerIdRequest{
		PassengerId: passengerId,
		Offset: offset,
		Limit: limit,
	})

	if err != nil {
		fmt.Println("could not get feedbacks by passenger:", err)
		return
	}

	fmt.Println(r.Msg)
	if r.ErrorCode == pb.Error_SUCCESS && len(r.Data) > 0 {
		for _, v := range r.Data {
			fmt.Println(v)
		}
	}
}

func getFeedbackByBookingCodeFromInput(c pb.FeedbackClient) {
	var bookingCode string
	fmt.Print("Input Booking Code: ")

	_, err := fmt.Scan(&bookingCode)

	if err != nil {
		log.Fatal(err)
		return
	}
	if len(strings.TrimSpace(bookingCode)) == 0 {
		fmt.Println("Booking code expected at least one character.")
		return
	}

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetFeedbackByBookingCode(ctx, &pb.PassengerFeedbackByBookingCodeRequest{
		BookingCode: bookingCode,
	})

	if err != nil {
		fmt.Println("could not get feedbacks by booking code:", err)
		return
	}

	fmt.Println(r.Msg)
	if r.ErrorCode == pb.Error_SUCCESS {
		fmt.Println(r.Data.String())
	}

}

func deleteFeedbackByPassengerIdFromInput(c pb.FeedbackClient) {
	var passengerId int32

	fmt.Print("Input Passenger Id: ")

	_, err := fmt.Scan(&passengerId)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DeleteFeedbackByPassengerId(ctx, &pb.DeletePassengerFeedbackByPassengerIdRequest{PassengerId:passengerId})
	if err != nil {
		fmt.Println("could not delete feedback by passenger:", err)
		return
	}

	fmt.Println(r.Msg)
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
	}
	defer conn.Close()
	c := pb.NewFeedbackClient(conn)

	fmt.Println("Before you choose a method, I'll add some passenger feedbacks to the local variable for passenger Id 1 " +
		"with booking from DRIVE01 to DRIVE010. ")
	number := 10

	chanInitFeedback := make(chan bool)

	go generateSomeFeedback(c, chanInitFeedback, number)

	<-chanInitFeedback

	fmt.Println("DONE")
	fmt.Println("Now is your turn. Following the instruction: ")

	for {
		fmt.Println("----Choose Method----")
		for i := 1; i <= len(methods); i++ {
			fmt.Println(i, ":", methods[i])
		}

		fmt.Print("> ")
		var methodID int
		_, err := fmt.Scan(&methodID)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if methodID < 1 || methodID > len(methods) {
			fmt.Println("Invalid Method")
			continue
		}

		fmt.Println("You choose:", methods[methodID])

		switch methodID {
		case 1:
			addPassengerFromInput(c)
			fmt.Println("Thank you for the information")
			continue
		case 2:
			getFeedbackByPassengerIdFromInput(c)
			continue
		case 3:
			getFeedbackByBookingCodeFromInput(c)
			continue
		case 4:
			deleteFeedbackByPassengerIdFromInput(c)
			continue
		case 5:
			fmt.Println("The program stops")
			return
		}
	}
}
