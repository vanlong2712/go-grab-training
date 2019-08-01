package main

import (
	pb "../proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
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

func (s *server) GetFeedbackByPassengerId(ctx context.Context, in *pb.GetPassengerFeedbackByPassengerIdRequest) (out *pb.PassengerFeedbackSliceResponse, err error) {
	out = new(pb.PassengerFeedbackSliceResponse)
	if len(feedbacks) > 0 {
		var sliceFeedbacks []string
		for k, v := range feedbacks {
			if v.PassengerId == in.PassengerId {
				sliceFeedbacks = append(sliceFeedbacks, k)
			}
		}

		if dataLength := len(sliceFeedbacks); dataLength < int(in.Offset) {
			sliceFeedbacks = []string{}
		} else if dataLength < int(in.Offset+in.Limit) {
			sliceFeedbacks = sliceFeedbacks[in.Offset:]
		} else {
			sliceFeedbacks = sliceFeedbacks[in.Offset:in.Limit]
		}

		if num := len(sliceFeedbacks); num > 0 {
			out.Msg = "The passenger has " + strconv.Itoa(num) + " feedbacks"
			for i := 0; i < num ; i++ {
				v := feedbacks[sliceFeedbacks[i]]
				out.Data = append(out.Data, &v)
			}

		} else {
			out.Msg = "The passenger has no feedback"
		}

		return out, nil
	}

	out.Msg = "The Passenger has no feedback"
	return out, nil
}

func (s *server) GetFeedbackByBookingCode(ctx context.Context, in *pb.PassengerFeedbackByBookingCodeRequest) (*pb.PassengerFeedbackResponse, error) {
	if feedback, ok := feedbacks[in.BookingCode]; ok {
		return &pb.PassengerFeedbackResponse{Data: &feedback, Msg: "Success!", ErrorCode: pb.Error_SUCCESS}, nil
	}
	return &pb.PassengerFeedbackResponse{Msg: "There is no feedback at booking code " + in.BookingCode, ErrorCode: pb.Error_FAIL}, nil
}

func (s *server) DeleteFeedbackByPassengerId(ctx context.Context, in *pb.DeletePassengerFeedbackByPassengerIdRequest) (*pb.ErrorCodeAndMessageResponse, error) {
	for k,v := range feedbacks {
		if v.PassengerId == in.PassengerId {
			delete(feedbacks, k)
		}
	}
	return &pb.ErrorCodeAndMessageResponse{Msg: "Success"}, nil
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
