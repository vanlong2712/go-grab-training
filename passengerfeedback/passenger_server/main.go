package main

import (
	"context"
	"log"
	"net"

	pb "github.com/vanlong2712/go-grab-training/passengerfeedback/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

var feedbacks = map[string]pb.PassengerFeedback{}

type server struct{}

func (s *server) AddPassengerFeedback(ctx context.Context, in *pb.PassengerFeedbackRequest) (*pb.PassengerFeedbackResponse, error) {
	if _, ok := feedbacks[in.BookingCode]; !ok {
		pfb := pb.PassengerFeedback{
			BookingCode: in.BookingCode,
			PassengerId: in.PassengerId,
			Feedback:    in.Feedback,
		}
		feedbacks[in.BookingCode] = pfb
		return &pb.PassengerFeedbackResponse{
			Data:      &pfb,
			Msg:       "Thank you for your feedback",
			ErrorCode: pb.Error_SUCCESS,
		}, nil
	}
	return &pb.PassengerFeedbackResponse{
		Msg:       "Feedback exists with the bookingCode",
		ErrorCode: pb.Error_FAIL,
	}, nil
}

func (s *server) GetFeedbackByPassengerId(ctx context.Context, in *pb.GetPassengerFeedbackByPassengerIdRequest) (*pb.PassengerFeedbackSpliceResponse, error) {
	panic("implement me")
}

func (s *server) GetFeedbackByBookingCode(ctx context.Context, in *pb.PassengerFeedbackByBookingCodeRequest) (*pb.PassengerFeedbackResponse, error) {
	panic("implement me")
}

func (s *server) DeleteFeedbackByPassengerId(ctx context.Context, in *pb.DeletePassengerFeedbackByPassengerIdRequest) (*pb.PassengerFeedbackResponse, error) {
	panic("implement me")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFeedbackServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
