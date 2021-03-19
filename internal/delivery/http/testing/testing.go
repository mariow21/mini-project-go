package testing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	teEntity "testing/internal/entity/testing"

	jaegerLog "testing/pkg/log"
	"testing/pkg/response"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// ITestingSvc is an interface to Testing Service
// Masukkan function dari service ke dalam interface ini
type ITestingSvc interface {
	GetTesting(ctx context.Context, _token string) error
	GetDataByName(ctx context.Context, nama string) ([]teEntity.Testing, error)
	GetDataByCity(ctx context.Context, kota string) ([]teEntity.Testing, error)
	GetDataByID(ctx context.Context, users []int) ([]teEntity.Testing, error)
	InsertDataArray(ctx context.Context, users []teEntity.Testing) error
	EditDataByID(ctx context.Context, users teEntity.Testing) error
}

type (
	// Handler ...
	Handler struct {
		testingSvc ITestingSvc
		tracer     opentracing.Tracer
		logger     jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(is ITestingSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		testingSvc: is,
		tracer:     tracer,
		logger:     logger,
	}
}

// TestingHandler will receive request and return response
func (h *Handler) TestingHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      context.Context
		_token   string
		resp     *response.Response
		result   interface{}
		metadata interface{}
		err      error
		errRes   response.Error
	)
	ctx = context.Background()
	_token = ""
	if r.Header.Get("Authorization") != "" {
		_token = r.Header.Get("Authorization")
		if !strings.Contains(_token, "Token ") {
			r.Method = "403"
		}
	}

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := h.tracer.StartSpan("TestingHandler", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	resp = &response.Response{}
	defer resp.RenderJSON(w, r)

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	if err == nil {
		switch r.Method {
		// Check if request method is GET
		case http.MethodGet:

			len := len(r.URL.Query())
			switch len {
			case 2:
				_, namaOK := r.URL.Query()["nama"]
				_, kotaOK := r.URL.Query()["kota"]
				if namaOK {
					result, err = h.testingSvc.GetDataByName(context.Background(), r.FormValue("nama"))
				} else if kotaOK {
					result, err = h.testingSvc.GetDataByCity(context.Background(), r.FormValue("kota"))
				} else {
					err = errors.New("400")
				}
			case 1:
				_, userIDOK := r.URL.Query()["user"]
				if userIDOK {
					var users teEntity.User
					body, _ := ioutil.ReadAll(r.Body)

					json.Unmarshal(body, &users)
					fmt.Println("Masuk disini", users)
					result, err = h.testingSvc.GetDataByID(context.Background(), users.UserID)
				}

			}
		// Check if request method is POST
		case http.MethodPost:
			_, postOK := r.URL.Query()["post"]
			body, _ := ioutil.ReadAll(r.Body)
			var JSON teEntity.User
			json.Unmarshal(body, &JSON)
			if postOK {	
				fmt.Println(JSON.UserInsert)
				err = h.testingSvc.InsertDataArray(context.Background(), JSON.UserInsert)
			}

		// Check if request method is PUT
		case http.MethodPut:
			var users teEntity.Testing
			body, _ := ioutil.ReadAll(r.Body)

			json.Unmarshal(body, &users)
			err = h.testingSvc.EditDataByID(context.Background(), users)
		// Check if request method is DELETE
		case http.MethodDelete:

		default:
			err = errors.New("400")
		}
	}

	// If anything from service or data return an error
	if err != nil {
		// Error response handling
		errRes = response.Error{
			Code:   101,
			Msg:    "101 - Data Not Found",
			Status: true,
		}
		// If service returns an error
		if strings.Contains(err.Error(), "service") {
			// Replace error with server error
			errRes = response.Error{
				Code:   500,
				Msg:    "500 - Internal Server Error",
				Status: true,
			}
		}
		// If error 401
		if strings.Contains(err.Error(), "401") {
			// Replace error with server error
			errRes = response.Error{
				Code:   401,
				Msg:    "401 - Unauthorized",
				Status: true,
			}
		}
		// If error 403
		if strings.Contains(err.Error(), "403") {
			// Replace error with server error
			errRes = response.Error{
				Code:   403,
				Msg:    "403 - Forbidden",
				Status: true,
			}
		}
		// If error 400
		if strings.Contains(err.Error(), "400") {
			// Replace error with server error
			errRes = response.Error{
				Code:   400,
				Msg:    "400 - Bad Request",
				Status: true,
			}
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		h.logger.For(ctx).Error("HTTP request error", zap.String("method", r.Method), zap.Stringer("url", r.URL), zap.Error(err))
		resp.StatusCode = errRes.Code
		resp.Error = errRes
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", r.Method), zap.Stringer("url", r.URL))
}
