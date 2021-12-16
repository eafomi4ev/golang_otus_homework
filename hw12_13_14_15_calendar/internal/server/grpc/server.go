package internalgrpc

//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=.

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/utils/idgen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	server  *grpc.Server
	logg    app.Logger
	app     *app.App
	address string
}

func NewServer(host string, port string, app *app.App, logger app.Logger) *Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggerInterceptor(logger)))
	grpc.ChainStreamInterceptor()
	server := &Server{
		server:  grpcServer,
		logg:    logger,
		app:     app,
		address: net.JoinHostPort(host, port),
	}

	RegisterCalendarServer(grpcServer, server)

	return server
}

func loggerInterceptor(logger app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		logger.Info(fmt.Sprintf("date: %s, method: %s, request: %+v", time.Now(), info.FullMethod, req))
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func (s Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen an address: " + err.Error())
	}

	go func() {
		s.logg.Info(fmt.Sprintf("grpc server is running on %s", s.address))
		if err := s.server.Serve(lsn); err != nil {
			s.logg.Error("failed to start grpc server: " + err.Error())
		}
		s.logg.Info("grpc servers has been stopped")
	}()

	go func() {
		<-ctx.Done()
		s.server.Stop()
	}()

	return nil
}

func (s Server) Stop() {
	s.server.GracefulStop()
}

func (s Server) Add(ctx context.Context, e *Event) (*Event, error) {
	eventID, _ := idgen.PrefixedID("EV")

	event := storage.Event{
		ID:          eventID,
		Title:       e.Title,
		EventDate:   e.EventDate.AsTime(),
		Duration:    time.Duration(e.Duration),
		Description: e.Description,
		UserID:      e.UserID,
		RemindIn:    time.Duration(e.RemindIn),
	}

	createdEvent, err := s.app.Storage.Add(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("grpc add event error: %w", err)
	}

	respEvent := &Event{
		Id:          createdEvent.ID,
		Title:       createdEvent.Title,
		EventDate:   timestamppb.New(createdEvent.EventDate),
		Duration:    int64(createdEvent.Duration),
		Description: createdEvent.Description,
		UserID:      createdEvent.UserID,
		RemindIn:    int64(createdEvent.RemindIn),
	}

	return respEvent, nil
}

func (s Server) Update(ctx context.Context, e *Event) (*Event, error) {
	event := storage.Event{
		ID:          e.Id,
		Title:       e.Title,
		EventDate:   e.EventDate.AsTime(),
		Duration:    time.Duration(e.Duration),
		Description: e.Description,
		UserID:      e.UserID,
		RemindIn:    time.Duration(e.RemindIn),
	}

	updatedEvent, err := s.app.Storage.Update(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("grpc update event error: %w", err)
	}

	respEvent := &Event{
		Id:          updatedEvent.ID,
		Title:       updatedEvent.Title,
		EventDate:   timestamppb.New(updatedEvent.EventDate),
		Duration:    int64(updatedEvent.Duration),
		Description: updatedEvent.Description,
		UserID:      updatedEvent.UserID,
		RemindIn:    int64(updatedEvent.RemindIn),
	}

	return respEvent, nil
}

func (s Server) Delete(ctx context.Context, request *DeleteRequest) (*emptypb.Empty, error) {
	err := s.app.Storage.Delete(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("grpc delete event error: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s Server) ListPerPeriod(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	return getEventList(ctx, request, s)
}

func (s Server) mustEmbedUnimplementedCalendarServer() {}

func storeEvToGRPCEv(sEv storage.Event) *Event {
	return &Event{
		Id:          sEv.ID,
		Title:       sEv.Title,
		EventDate:   timestamppb.New(sEv.EventDate),
		Duration:    int64(sEv.Duration),
		Description: sEv.Description,
		UserID:      sEv.UserID,
		RemindIn:    int64(sEv.RemindIn),
	}
}

func getEventList(ctx context.Context, request *ListRequest, s Server) (*ListResponse, error) {
	var events []storage.Event
	var err error
	switch request.PeriodName {
	case "day":
		events, err = s.app.Storage.ListPerDay(ctx, request.StartDay.AsTime())
	case "week":
		events, err = s.app.Storage.ListPerWeek(ctx, request.StartDay.AsTime())
	case "month":
		events, err = s.app.Storage.ListPerMonth(ctx, request.StartDay.AsTime())
	}

	if err != nil {
		return nil, fmt.Errorf("grpc getEventList event error: %w", err)
	}

	respEvents := make([]*Event, 0)
	for _, e := range events {
		respEvents = append(respEvents, storeEvToGRPCEv(e))
	}

	resp := ListResponse{
		Events: respEvents,
		Error:  nil,
	}

	return &resp, nil
}
