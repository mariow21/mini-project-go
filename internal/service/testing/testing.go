package testing

import (
	"context"
	"testing/internal/entity/auth"
	teEntity "testing/internal/entity/testing"
	"testing/pkg/errors"
	jaegerLog "testing/pkg/log"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	GetDataByName(ctx context.Context, nama string) ([]teEntity.Testing, error)
	GetDataByCity(ctx context.Context, kota string) ([]teEntity.Testing, error)
	GetDataByID(ctx context.Context, users int) (teEntity.Testing, error)
	InsertDataArray(ctx context.Context, users teEntity.Testing) error
	EditDataByID(ctx context.Context, users teEntity.Testing) error
}

// AuthData ...
type AuthData interface {
	CheckAuth(ctx context.Context, _token, code string) (auth.Auth, error)
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	data     Data
	authData AuthData
	tracer   opentracing.Tracer
	logger   jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(data Data, authData AuthData, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		data:     data,
		authData: authData,
		tracer:   tracer,
		logger:   logger,
	}
}

// GetTesting ...
func (s Service) GetTesting(ctx context.Context, _token string) error {
	// Check if have span on context
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("GetTesting", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	// CHECK AUTH
	s.logger.For(ctx).Info("Check Auth", zap.String("code", "191"))
	auth, err := s.authData.CheckAuth(ctx, _token, "191")
	if err != nil {
		s.logger.For(ctx).Error("Failed to check auth", zap.Error(err))
		return errors.Wrap(err, "[SERVICE][GetTesting]")
	}
	if auth.Error.Status == true {
		s.logger.For(ctx).Error("401 Unauthorized")
		return errors.Wrap(errors.New("401 Unauthorized"), "[SERVICE][GetTesting]")
	}
	s.logger.For(ctx).Info("End Check Auth")
	// END CHECK AUTH

	return nil
}

func (s Service) GetDataByName(ctx context.Context, nama string) ([]teEntity.Testing, error) {
	var (
		users []teEntity.Testing
		err   error
	)

	users, err = s.data.GetDataByName(ctx, nama)
	if err != nil {
		return users, errors.Wrap(err, "[SERVICE][GetDataByName]")
	}
	return users, err
}

func (s Service) GetDataByCity(ctx context.Context, kota string) ([]teEntity.Testing, error) {
	var (
		users []teEntity.Testing
		err   error
	)

	users, err = s.data.GetDataByCity(ctx, kota)
	if err != nil {
		return users, errors.Wrap(err, "[SERVICE][GetDataByCity]")
	}
	return users, err
}

func (s Service) GetDataByID(ctx context.Context, users []int) ([]teEntity.Testing, error) {
	var (
		result  teEntity.Testing
		results []teEntity.Testing
		err     error
	)

	for _, user := range users {
		result, err = s.data.GetDataByID(ctx, user)
		if err != nil {
			return results, errors.Wrap(err, "[SERVICE][GetDataByID]")
		}
		results = append(results, result)
	}

	return results, err
}

func (s Service) InsertDataArray(ctx context.Context, users []teEntity.Testing) error {
	var (
		err error
	)

	for _, user := range users {
		err = s.data.InsertDataArray(ctx, user)

		if err != nil {
			return errors.Wrap(err, "[SERVICE][InsertDataArray]")
		}
	}

	return err
}

func (s Service) EditDataByID(ctx context.Context, users teEntity.Testing) error {
	var err error

	err = s.data.EditDataByID(ctx, users)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][EditDataByID]")
	}

	return err
}
