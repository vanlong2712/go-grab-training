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

var feedbacks = map[int32][]pb.PassengerFeedback{}

type server struct{}

func (s *server) AddPassengerFeedback(ctx context.Context, in *pb.PassengerFeedbackRequest) (*pb.PassengerFeedbackResponse, error) {
	if feedbackByPassenger, ok := feedbacks[in.PassengerId]; ok {
		for _, fb := range feedbackByPassenger {
			if fb.BookingCode == in.BookingCode {
				return &pb.PassengerFeedbackResponse{
					Msg:       "Feedback exists",
					ErrorCode: pb.Error_FAIL,
				}, nil
			}
		}
	}

	pfb := pb.PassengerFeedback{
		BookingCode: in.BookingCode,
		PassengerId: in.PassengerId,
		Feedback:    in.Feedback,
	}
	feedbacks[in.PassengerId] = append(feedbacks[in.PassengerId], pfb)
	return &pb.PassengerFeedbackResponse{
		Data:      &pfb,
		Msg:       "Thank you for your feedback",
		ErrorCode: pb.Error_SUCCESS,
	}, nil
}

func (s *server) GetFeedbackByPassengerId(ctx context.Context, in *pb.GetPassengerFeedbackByPassengerIdRequest) (out *pb.PassengerFeedbackSliceResponse, err error) {
	out = new(pb.PassengerFeedbackSliceResponse)
	if len(feedbacks) > 0 {
		var sliceFeedbacks []int32
		passengerFeedbacks := feedbacks[in.PassengerId]

		if dataLength := int32(len(passengerFeedbacks)); dataLength < in.Offset {
			sliceFeedbacks = []int32{}
		} else if dataLength < in.Offset+in.Limit {
			for i := in.Offset; i < dataLength; i++ {
				sliceFeedbacks = append(sliceFeedbacks, i)
			}
		} else {
			for i := in.Offset; i < in.Limit; i++ {
				sliceFeedbacks = append(sliceFeedbacks, i)
			}
		}

		if num := len(sliceFeedbacks); num > 0 {
			out.Msg = "The passenger has " + strconv.Itoa(num) + " feedbacks"
			for i := 0; i < num ; i++ {
				v := feedbacks[in.PassengerId][i]
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

func (s *server) GetFeedbackByBookingCode(ctx context.Context, in *pb.PassengerFeedbackByBookingCodeRequest) (*pb.PassengerFeedbackSliceResponse, error) {
	out := new(pb.PassengerFeedbackSliceResponse)
	if len(feedbacks) > 0 {
		for _, feedbackByPassenger := range feedbacks {
			for k, fb := range feedbackByPassenger {
				if fb.BookingCode == in.BookingCode {
					out.Data = append(out.Data, &feedbackByPassenger[k])
				}
			}
		}
		if dataLength := int32(len(out.Data)); dataLength < in.Offset {
			out.Data = out.Data[0:0]
		} else if dataLength < in.Offset+in.Limit {
			out.Data = out.Data[in.Offset:]
		} else {
			out.Data = out.Data[in.Offset:in.Limit]
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
		if _, ok := feedbacks[in.PassengerId]; ok {
			delete(feedbacks, in.PassengerId)
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
