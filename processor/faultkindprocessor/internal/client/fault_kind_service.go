package client

import (
	"context"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/pkg/spangroup"
	"github.com/teanoon/opentelemetry-collector-contrib/processor/faultkindprocessor/ent/schema"
	"go.uber.org/zap"
)

type FaultKind int

const (
	BusinessFault FaultKind = iota
	SystemFault
)

func fromStringValue(value string) FaultKind {
	switch value {
	case "BusinessFault":
		return BusinessFault
	default:
		return SystemFault
	}
}

func (f FaultKind) String() string {
	switch f {
	case BusinessFault:
		return "BusinessFault"
	default:
		return "SystemFault"
	}
}

type FaultKindService interface {
	Start(ctx context.Context)
	MatchFaultKind(ctx context.Context, resourceAttributes *map[string]any, spanAttributes *map[string]any) string
	Shutdown(ctx context.Context) error
}

type FaultKindServiceImpl struct {
	logger          *zap.Logger
	repository      FaultKindRepository
	cacheTtlMinutes time.Duration
	groups          *spangroup.SpanGroups

	ticker *time.Ticker
}

func CreateFaultKindServiceImpl(logger *zap.Logger, repository FaultKindRepository, cacheTtlMinutes time.Duration) *FaultKindServiceImpl {
	return &FaultKindServiceImpl{logger: logger, repository: repository, cacheTtlMinutes: cacheTtlMinutes}
}

func (s *FaultKindServiceImpl) Start(ctx context.Context) {
	s.ticker = time.NewTicker(s.cacheTtlMinutes * time.Minute)

	go func(ctx context.Context) {
		for ; ; <-s.ticker.C {
			s.logger.Info("Building cache")
			err := s.buildCache(ctx)
			if err != nil {
				s.logger.Sugar().Errorf("Error when building cache: %s\n", err)
			}
		}
	}(ctx)
}

func (s *FaultKindServiceImpl) buildCache(ctx context.Context) error {
	if s.repository == nil {
		return nil
	}

	definitions, err := s.repository.FindFaultKindDefinitions(ctx)
	if err != nil {
		return err
	}

	data := make(map[*spangroup.SpanGroupDefinitions]string, 2)
	build := func(definitions []*schema.FaultKindCondition, faultKind FaultKind) {
		group := spangroup.SpanGroupDefinitions{}
		for _, definition := range definitions {
			group = append(group, spangroup.SpanGroupDefinition{
				Column: definition.Column,
				Op:     definition.Op,
				Value:  spangroup.CreateDefinitionValue(definition.Value),
			})
		}
		data[&group] = faultKind.String()
	}
	build(definitions.Business, BusinessFault)
	build(definitions.System, SystemFault)
	s.groups = spangroup.CreateSpanGroup(data)

	return nil
}

func (s *FaultKindServiceImpl) MatchFaultKind(ctx context.Context, resourceAttributes *map[string]any, spanAttributes *map[string]any) string {
	attributes := make(map[string]interface{}, len(*resourceAttributes)+len(*spanAttributes))
	for key, value := range *resourceAttributes {
		attributes[key] = value
	}
	for key, value := range *spanAttributes {
		attributes[key] = value
	}
	kind := s.groups.Get(&attributes)
	if len(kind) == 0 {
		return ""
	}
	return fromStringValue(kind[0]).String()
}

func (s *FaultKindServiceImpl) Shutdown(ctx context.Context) error {
	s.ticker.Stop()
	return s.repository.Shutdown(ctx)
}
