package ts

import (
	"context"
	"fmt"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

type Storage interface {
	CreateSchema(name string, field []api.Field, defaultSort string) error
	Index(ctx context.Context, input IndexInput) error
	RetrieveById(ctx context.Context, input RetrieveInput) (*RetrieveOutput, error)
	Search(ctx context.Context, input SearchInput) (*SearchOutput, error)
}

type StorageConfig struct {
	TSClient *typesense.Client
}

func (c StorageConfig) Validate() error {
	if c.TSClient == nil {
		return fmt.Errorf("missing `TypeSense Client`")
	}
	return nil
}

func NewStorage(cfg StorageConfig) (Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &storage{tsClient: cfg.TSClient}
	return s, nil
}

type storage struct {
	tsClient *typesense.Client
}

func (s *storage) CreateSchema(name string, field []api.Field, defaultSort string) error {
	schema := &api.CollectionSchema{
		Name:                name,
		Fields:              field,
		DefaultSortingField: &defaultSort,
	}

	_, err := s.tsClient.Collections().Create(schema)
	if err != nil {
		return fmt.Errorf("unable to create schema due: %w", err)
	}

	return nil
}

func (s *storage) Index(ctx context.Context, input IndexInput) error {
	_, err := s.tsClient.Collection(input.CollectionName).Documents().Create(input.Doc)
	if err != nil {
		return fmt.Errorf("unable to index a document due: %w", err)
	}
	return nil
}

func (s *storage) RetrieveById(ctx context.Context, input RetrieveInput) (*RetrieveOutput, error) {
	res, err := s.tsClient.Collection(input.CollectionName).Document(input.ID).Retrieve()
	if err != nil {
		return nil, fmt.Errorf("unale to retrieve a document due: %w", err)
	}

	output := RetrieveOutput{
		CollectionName: input.CollectionName,
		Doc: &MapDoc{
			Source: res,
		},
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse retrieve result due: %w", err)
	}
	return &output, nil
}

func (s *storage) Search(ctx context.Context, input SearchInput) (*SearchOutput, error) {
	res, err := s.tsClient.Collection(input.CollectionName).Documents().Search(&api.SearchCollectionParams{
		Q:       input.Query,
		QueryBy: input.QueryBy,
		Page:    pointer.Int(input.Page),
		PerPage: pointer.Int(input.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to execute search due: %w", err)
	}

	var output SearchOutput
	output.Found = *res.Found
	for _, hit := range *res.Hits {
		output.Hits = append(output.Hits, &MapDoc{
			Source: *hit.Document,
		})
	}
	return &output, nil
}
