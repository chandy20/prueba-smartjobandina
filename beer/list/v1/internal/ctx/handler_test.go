package ctx

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

//beerRepositoryMock struct to simulate beer repository behavior
type beerRepositoryMock struct {
	mock.Mock
}

func (b *beerRepositoryMock) List() ([]model.Beer, error) {
	args := b.Called()
	return args.Get(0).([]model.Beer), args.Error(1)
}

//go:embed golden_files/successResponse.json
var successResponse []byte

func TestHandler_Handler(t *testing.T) {
	beersToReturn := []model.Beer{
		{
			ID:       1,
			Name:     "Pilsen",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2400,
			Currency: "COP",
		},
		{
			ID:       2,
			Name:     "Brava",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2000,
			Currency: "COP",
		},
		{
			ID:       3,
			Name:     "Corona",
			Brewery:  "Bavaria",
			Country:  "Mexico",
			Price:    200,
			Currency: "MXN",
		},
		{
			ID:       4,
			Name:     "Budweiser",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2.5,
			Currency: "USD",
		},
		{
			ID:       5,
			Name:     "Leona",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2000,
			Currency: "COP",
		},
		{
			ID:       6,
			Name:     "Reds",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    3,
			Currency: "EUR",
		},
	}
	headers := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}

	type mocks struct {
		beersRepository *beerRepositoryMock
	}
	type fields struct {
		logger *logrus.Logger
	}
	type args struct {
		ctx context.Context
		req events.APIGatewayProxyRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		mocks   mocks
		mocker  func(m mocks)
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "should_return_error_listing_beers",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("List").Return([]model.Beer{}, errors.New("error")).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
				Body:       `{"message":"error"}`,
			},
			wantErr: false,
		},
		{
			name: "should_return_an_empty_response_because_there_are_no_beers",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("List").Return([]model.Beer{}, nil).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusAccepted,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
			wantErr: false,
		},
		{
			name: "should_return_an_success_response",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("List").Return(beersToReturn, nil).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    headers,
				Body:       string(successResponse),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocker(tt.mocks)
			h := NewHandler(tt.mocks.beersRepository, tt.fields.logger)
			got, err := h.Handler(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler() got = %v, want %v", got, tt.want)
			}
			tt.mocks.beersRepository.AssertExpectations(t)
		})
	}
}
