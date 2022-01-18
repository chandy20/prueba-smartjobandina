package ctx

import (
	"context"
	_ "embed"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"net/http"
	"reflect"
	"testing"
)

//beersRepositoryMock mock for represent beers repository
type beersRepositoryMock struct {
	mock.Mock
}

func (b *beersRepositoryMock) Find(ID int) (model.Beer, error) {
	args := b.Called(ID)
	return args.Get(0).(model.Beer), args.Error(1)
}

//go:embed golden_files/successResponse.json
var successResponse []byte

func TestHandler_Handler(t *testing.T) {
	headers := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}

	type mocks struct {
		beersRepository *beersRepositoryMock
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
		mocks   mocks
		mocker  func(m mocks)
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "should_return_error_because_beer_id_is_not_present",
			fields: fields{
				logger: logrus.New(),
			},
			mocks: mocks{
				beersRepository: &beersRepositoryMock{},
			},
			mocker: func(m mocks) {},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         headers,
				Body:            `{"message":"beerID_can_not_be_empty"}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "should_return_error_because_beer_id_is_not_a_number",
			fields: fields{
				logger: logrus.New(),
			},
			mocks: mocks{
				beersRepository: &beersRepositoryMock{},
			},
			mocker: func(m mocks) {},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{
						"beerID": "a",
					},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         headers,
				Body:            `{"message":"beerID_is_not_a_number"}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "should_return_error_finding_beer",
			fields: fields{
				logger: logrus.New(),
			},
			mocks: mocks{
				beersRepository: &beersRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", 1).Return(model.Beer{}, errors.New("error")).Once()
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{
						"beerID": "1",
					},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				Headers:         headers,
				Body:            `{"message":"error"}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "should_return_error_because_beer_id_does_not_exist",
			fields: fields{
				logger: logrus.New(),
			},
			mocks: mocks{
				beersRepository: &beersRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", 1).Return(model.Beer{}, nil).Once()
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{
						"beerID": "1",
					},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusNotFound,
				Headers:         headers,
				Body:            `{"message":"beerID_does_not_exist"}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "should_return_a_success_response",
			fields: fields{
				logger: logrus.New(),
			},
			mocks: mocks{
				beersRepository: &beersRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", 1).Return(
					model.Beer{
						ID:       1,
						Name:     "Pilsen",
						Brewery:  "Bavaria",
						Country:  "Colombia",
						Price:    2400,
						Currency: "COP",
					}, nil).Once()
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{
						"beerID": "1",
					},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         headers,
				Body:            string(successResponse),
				IsBase64Encoded: false,
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
		})
	}
}
