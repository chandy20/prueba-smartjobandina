package ctx

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

//beerRepositoryMock struct to simulate beer repository behavior
type beerRepositoryMock struct {
	mock.Mock
}

func (b *beerRepositoryMock) Find(id int) (model.Beer, error) {
	args := b.Called(id)
	return args.Get(0).(model.Beer), args.Error(1)
}

func (b *beerRepositoryMock) Save(beer model.Beer) error {
	return b.Called(beer).Error(0)
}

//go:embed  golden_files/success_message.json
var successMessage []byte

//go:embed golden_files/error_message.json
var errorMessage []byte

func TestHandler_Handler(t *testing.T) {
	headers := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}
	var beer model.Beer
	err := json.Unmarshal(successMessage, &beer)
	if err != nil {
		t.Error(err)
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
			name: "should_return_error_because_request_is_invalid",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
					Body:    string(errorMessage),
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    headers,
				Body:       `{"message":"json: cannot unmarshal string into Go struct field Beer.id of type int"}`,
			},
			wantErr: false,
		},
		{
			name: "should_return_error_because_beer_is_already_created",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
					Body:    string(successMessage),
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", beer.ID).Return(beer, nil).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusConflict,
				Headers:    headers,
				Body:       `{"message":"error_beer_already_created"}`,
			},
			wantErr: false,
		},
		{
			name: "should_return_finding_beer",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
					Body:    string(successMessage),
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", beer.ID).Return(model.Beer{}, errors.New("error")).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
				Body:       `{"message":"error"}`,
			},
			wantErr: false,
		},
		{
			name: "should_return_error_saving_beer",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
					Body:    string(successMessage),
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", beer.ID).Return(model.Beer{}, nil).Once()
				m.beersRepository.On("Save", beer).Return(errors.New("error")).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
				Body:       `{"message":"error"}`,
			},
			wantErr: false,
		},
		{
			name: "should_return_success_response",
			fields: fields{
				logger: logrus.New(),
			},
			args: args{
				ctx: context.Background(),
				req: events.APIGatewayProxyRequest{
					Headers: headers,
					Body:    string(successMessage),
				},
			},
			mocks: mocks{
				beersRepository: &beerRepositoryMock{},
			},
			mocker: func(m mocks) {
				m.beersRepository.On("Find", beer.ID).Return(model.Beer{}, nil).Once()
				m.beersRepository.On("Save", beer).Return(nil).Once()
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusCreated,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: "",
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
