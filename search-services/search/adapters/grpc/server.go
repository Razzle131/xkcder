package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/core"
)

type Server struct {
	searchpb.UnimplementedSearchServer
	service core.Searcher
}

func NewServer(service core.Searcher) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Search(ctx context.Context, in *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	resp, err := s.service.Search(ctx, int(in.Limit), in.Phrase)
	if err != nil {
		return nil, err
	}

	res := searchpb.SearchReply{
		Comics: make([]*searchpb.Comics, 0, len(resp)),
	}

	for _, c := range resp {
		res.Comics = append(res.Comics, &searchpb.Comics{
			Id:  int64(c.ID),
			Url: c.URL,
		})
	}

	return &res, nil
}

func (s *Server) SearchIndex(ctx context.Context, in *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	resp, err := s.service.SearchIndex(ctx, int(in.Limit), in.Phrase)
	if err != nil {
		return nil, err
	}

	res := searchpb.SearchReply{
		Comics: make([]*searchpb.Comics, 0, len(resp)),
	}

	for _, c := range resp {
		res.Comics = append(res.Comics, &searchpb.Comics{
			Id:  int64(c.ID),
			Url: c.URL,
		})
	}

	return &res, nil
}
