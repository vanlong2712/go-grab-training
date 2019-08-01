package main

import (
	pb "../proto"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	address = "localhost:50051"
)

var client pb.FeedbackClient

func handleAddFeedback(c *gin.Context) {
	var argument struct {
		PassengerId string
		BookingCode string
		Feedback    string
	}
	err := c.BindJSON(&argument)
	if err != nil {
		c.String(400, "invalid param \n")
		log.Print(err)
		return
	}
	if len(strings.TrimSpace(argument.PassengerId)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing passengerId"})
		return
	}
	passengerId, err := strconv.Atoi(argument.PassengerId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "PassengerId must be a number"})
		return
	}
	if len(strings.TrimSpace(argument.BookingCode)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing bookingCode"})
		return
	}
	if len(strings.TrimSpace(argument.Feedback)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing feedback"})
		return
	}
	res, err := addPassengerFeedback(argument.BookingCode, int32(passengerId), argument.Feedback)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Something happens"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_SUCCESS, "message": res.Msg, "data": res.Data})
}

func handleGetFeedbackByPassengerId(c *gin.Context) {
	strPassenger, b := c.GetQuery("passengerId")
	if !b || len(strings.TrimSpace(strPassenger)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing passengerId"})
		return
	}
	passengerId, err := strconv.Atoi(strPassenger)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "PassengerId must be a number"})
		return
	}
	strOffset, b := c.GetQuery("offset")
	if !b || len(strings.TrimSpace(strOffset)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing Offset"})
		return
	}
	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Offset must be a number"})
		return
	}
	strLimit, b := c.GetQuery("limit")
	if !b || len(strings.TrimSpace(strLimit)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing Limit"})
		return
	}
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Limit must be a number"})
		return
	}

	response, err := getFeedbackByPassengerId(int32(passengerId), int32(offset), int32(limit))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorCode": pb.Error_FAIL, "message": err.Error()})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"errorCode": response.ErrorCode, "message": response.Msg, "data": response.Data})
}

func handleGetFeedbackByBookingCode(c *gin.Context) {
	strBookingCode, b := c.GetQuery("bookingCode")
	if !b || len(strings.TrimSpace(strBookingCode)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing bookingCode"})
		return
	}
	strOffset, b := c.GetQuery("offset")
	if !b || len(strings.TrimSpace(strOffset)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing Offset"})
		return
	}
	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Offset must be a number"})
		return
	}
	strLimit, b := c.GetQuery("limit")
	if !b || len(strings.TrimSpace(strLimit)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing Limit"})
		return
	}
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Limit must be a number"})
		return
	}

	response, err := getFeedbackByBookingCode(strBookingCode, int32(offset), int32(limit))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorCode": pb.Error_FAIL, "message": err.Error()})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"errorCode": response.ErrorCode, "message": response.Msg, "data": response.Data})
}

func handleDeleteFeedbackByPassengerId(c *gin.Context) {
	var argument struct {
		PassengerId string
	}
	err := c.BindJSON(&argument)
	if err != nil {
		c.String(400, "invalid param \n")
		log.Print(err)
		return
	}
	if len(strings.TrimSpace(argument.PassengerId)) == 0 {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "Missing passengerId"})
		return
	}
	passengerId, err := strconv.Atoi(argument.PassengerId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errorCode": pb.Error_FAIL, "message": "PassengerId must be a number"})
		return
	}

	response, err := deleteFeedbackByPassengerId(int32(passengerId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorCode": pb.Error_FAIL, "message": err.Error()})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"errorCode": response.ErrorCode, "message": response.Msg})
}

func addPassengerFeedback(bookingCode string, passengerId int32, feedback string) (res *pb.PassengerFeedbackResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err = client.AddPassengerFeedback(ctx, &pb.PassengerFeedbackRequest{
		BookingCode: bookingCode,
		PassengerId: passengerId,
		Feedback:    feedback,
	})
	return
}

func getFeedbackByPassengerId(passengerId, offset, limit int32) (*pb.PassengerFeedbackSliceResponse, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := client.GetFeedbackByPassengerId(ctx, &pb.GetPassengerFeedbackByPassengerIdRequest{
		PassengerId: passengerId,
		Offset:      offset,
		Limit:       limit,
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

func getFeedbackByBookingCode(bookingCode string, offset, limit int32) (*pb.PassengerFeedbackSliceResponse, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := client.GetFeedbackByBookingCode(ctx, &pb.PassengerFeedbackByBookingCodeRequest{
		BookingCode: bookingCode,
		Offset:      offset,
		Limit:       limit,
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

func deleteFeedbackByPassengerId(passengerId int32)  (*pb.ErrorCodeAndMessageResponse, error)  {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := client.DeleteFeedbackByPassengerId(ctx, &pb.DeletePassengerFeedbackByPassengerIdRequest{PassengerId: passengerId})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
	}
	defer conn.Close()
	client = pb.NewFeedbackClient(conn)

	router := gin.Default()

	router.POST("/addFeedback", handleAddFeedback)
	//localhost:8080/addFeedback
	//body raw: {"bookingCode": "DRIVE01", "passengerId": "1", "feedback": "Good"}
	router.GET("/getFeedbackByPassengerId", handleGetFeedbackByPassengerId)
	//localhost:8080/getFeedbackByPassengerId?passengerId=1&offset=0&limit=10
	router.GET("/getFeedbackByBookingCode", handleGetFeedbackByBookingCode)
	//localhost:8080/getFeedbackByBookingCode?bookingCode=DRIVE02&offset=0&limit=10
	router.DELETE("/deleteFeedback", handleDeleteFeedbackByPassengerId)
	//localhost:8080/deleteFeedback
	//body raw: {"passengerId":"2"}
	router.Run()


}
