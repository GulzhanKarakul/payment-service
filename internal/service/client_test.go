package service_test

import (
	"context"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/service"
	"github.com/GulzhanKarakul/payment-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newClientSvc(t *testing.T) (service.ClientService, *mocks.MockClientRepository) {
	t.Helper()
	repo := mocks.NewMockClientRepository(t)
	svc := service.NewClientService(repo, testLogger())
	return svc, repo
}

// Create - ClientService
func TestClientService_Create_Success(t *testing.T) {
	svc, repo := newClientSvc(t)

	repo.EXPECT().GetByPhone(mock.Anything, "+77771156580").
		Return(domain.Client{}, domain.ErrClientNotFound).Once()
	repo.EXPECT().Create(mock.Anything, "+77771156580", "Gulzhan").
		Return(activeClient(), nil).Once()

	client, err := svc.Create(context.Background(), "+77771156580", "Gulzhan")

	require.NoError(t, err)
	assert.Equal(t, activeClient().ID, client.ID)
	assert.Equal(t, "+77771156580", client.Phone)
	assert.True(t, client.IsActive)
}

func TestClientService_Create_PhoneTaken(t *testing.T) {
	svc, repo := newClientSvc(t)

	repo.EXPECT().GetByPhone(mock.Anything, "+77771156580").Return(activeClient(), nil).Once()

	_, err := svc.Create(context.Background(), "+77771156580", "Gulzhan")

	require.ErrorIs(t, err, domain.ErrClientAlreadyExist)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestClientService_Create_RepoSystemError(t *testing.T) {
	svc, repo := newClientSvc(t)

	repo.EXPECT().GetByPhone(mock.Anything, "+77771156580").
		Return(domain.Client{}, ErrConnectionRefused).Once()

	_, err := svc.Create(context.Background(), "+77771156580", "Gulzhan")

	require.ErrorIs(t, err, ErrConnectionRefused)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

// GetByID
func TestClientService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct{
		name string
		id string
		setup func(repo *mocks.MockClientRepository)
		wantErr error
	} {
		{
			name: "success",
			id: activeClient().ID,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeClient().ID).
					Return(activeClient(), nil).Once()
			},
		},
		{
			name: "client not found",
			id: activeClient().ID,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeClient().ID).
					Return(domain.Client{}, domain.ErrClientNotFound).Once()
			},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "propogates system error",
			id: activeClient().ID,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeClient().ID).
					Return(domain.Client{}, ErrConnectionRefused).Once()
			},
			wantErr: ErrConnectionRefused,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newClientSvc(t)
			tt.setup(repo)

			client, err := svc.GetByID(ctx, tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, client.ID)
			assert.Equal(t, activeClient().Phone, client.Phone)
			assert.True(t, client.IsActive)
		})
	}
}

// GetByPhone
func TestClientService_GetByPhone(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		phone string
		setup func(repo *mocks.MockClientRepository)
		wantErr error
		wantActive *bool
	}{
		{
			name: "success - returns active client",
			phone: activeClient().Phone,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByPhone(mock.Anything, activeClient().Phone).
					Return(activeClient(), nil).Once()
			},
		},
		{
			name: "client not found",
			phone: activeClient().Phone,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByPhone(mock.Anything, activeClient().Phone).
					Return(domain.Client{}, domain.ErrClientNotFound).Once()
			},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "returns inactive client",
			phone: inactiveClient().Phone,
			setup: func(repo *mocks.MockClientRepository) {
				repo.EXPECT().GetByPhone(mock.Anything, inactiveClient().Phone).
					Return(inactiveClient(), nil).Once()
			},
			wantActive: func () *bool {b:= false; return &b}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newClientSvc(t)
			tt.setup(repo)

			client, err := svc.GetByPhone(ctx, tt.phone)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, client.ID)
			assert.Equal(t, tt.phone, client.Phone)
			if tt.wantActive != nil {
				assert.Equal(t, *tt.wantActive, client.IsActive)
			}
		})
	}
}