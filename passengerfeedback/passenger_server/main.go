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
	port = ":50051"
)

var feedbacks []pb.PassengerFeedback

type server struct{}

func (s *server) AddPassengerFeedback(ctx context.Context, in *pb.PassengerFeedbackRequest) (*pb.PassengerFeedbackResponse, error) {
		for _, fb := range feedbacks {
			if fb.BookingCode == in.BookingCode && fb.PassengerId == in.PassengerId {
				return &pb.PassengerFeedbackResponse{
					Msg:       "Feedback exists",
					ErrorCode: pb.Error_FAIL,
				}, nil
			}
		}


	pfb := pb.PassengerFeedback{
		BookingCode: in.BookingCode,
		PassengerId: in.PassengerId,
		Feedback:    in.Feedback,
	}
	feedbacks = append(feedbacks, pfb)
	return &pb.PassengerFeedbackResponse{
		Data:      &pfb,
		Msg:       "Thank you for your feedback",
		ErrorCode: pb.Error_SUCCESS,
	}, nil
}

func (s *server) GetFeedbackByPassengerId(ctx context.Context, in *pb.GetPassengerFeedbackByPassengerIdRequest) (out *pb.PassengerFeedbackSliceResponse, err error) {
	out = new(pb.PassengerFeedbackSliceResponse)
	if len(feedbacks) > 0 {
		var fbByPassengerId []pb.PassengerFeedback
		for i := 0; i < len(feedbacks); i++ {
			if in.PassengerId == feedbacks[i].PassengerId {
				fbByPassengerId = append(fbByPassengerId, feedbacks[i])
			}
		}

		for i := in.Offset; i < in.Offset + in.Limit; i++ {
			if len(fbByPassengerId) > int(i) {
				out.Data = append(out.Data, &fbByPassengerId[i])
			} else {
				break
			}
		}

		if num := len(out.Data); num > 0 {
			out.Msg = "The passenger has " + strconv.Itoa(num) + " feedbacks"
		} else {
			out.Msg = "The passenger has no feedback"
		}

		return out, nil
	}

	out.Msg = "The Passenger has no feedback"
	return out, nil
}

func (s *server) GetFeedbackByBookingCode(ctx context.Context, in *pb.PassengerFeedbackByBookingCodeRequest) (*pb.PassengerFeedbackSliceResponse, error) {
	out := new(pb.PassengerFeedbackSliceResponse)
	if len(feedbacks) > 0 {
		var fbByBookingCode []pb.PassengerFeedback
		for i := 0; i < len(feedbacks); i++ {
			if in.BookingCode == feedbacks[i].BookingCode {
				fbByBookingCode = append(fbByBookingCode, feedbacks[i])
			}
		}

		for i := in.Offset; i < in.Offset + in.Limit; i++ {
			if len(fbByBookingCode) > int(i) {
				out.Data = append(out.Data, &fbByBookingCode[i])
			} else {
				break
			}
		}

		if num := len(out.Data); num > 0 {
			out.Msg = "The booking code has " + strconv.Itoa(num) + " feedbacks"

		} else {
			out.Msg = "The booking code has no feedback"
		}

		return out, nil
	}
	out.Msg = "The booking code has no feedback"
	return out, nil
}

func (s *server) DeleteFeedbackByPassengerId(ctx context.Context, in *pb.DeletePassengerFeedbackByPassengerIdRequest) (*pb.ErrorCodeAndMessageResponse, error) {
	if len(feedbacks) > 0 {
		for i := 0; i < len(feedbacks); i++ {
			if feedbacks[i].PassengerId == in.PassengerId {
				feedbacks = append(feedbacks[:i], feedbacks[i+1:]...)
				i--
			}
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
