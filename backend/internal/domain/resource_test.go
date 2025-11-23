// backend/internal/domain/resource_test.go
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
)

func TestResource_Validate(t *testing.T) {
	capacity := 10
	tests := []struct {
		name     string
		resource domain.Resource
		wantErr  bool
	}{
		{
			name: "Valid meeting room",
			resource: domain.Resource{
				Name:     "Room A",
				Type:     domain.ResourceTypeMeetingRoom,
				Capacity: &capacity,
			},
			wantErr: false,
		},
		{
			name: "Valid equipment",
			resource: domain.Resource{
				Name: "Projector",
				Type: domain.ResourceTypeEquipment,
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			resource: domain.Resource{
				Type: domain.ResourceTypeMeetingRoom,
			},
			wantErr: true,
		},
		{
			name: "Missing type",
			resource: domain.Resource{
				Name: "Unknown",
			},
			wantErr: true,
		},
		{
			name: "Meeting room without capacity",
			resource: domain.Resource{
				Name: "Room B",
				Type: domain.ResourceTypeMeetingRoom,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resource.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResource_CanBeReservedBy(t *testing.T) {
	roleManager := domain.RoleManager
	resource := domain.Resource{
		IsActive:     true,
		RequiredRole: &roleManager,
	}
	inactiveResource := domain.Resource{
		IsActive: false,
	}

	manager := domain.User{Role: domain.RoleManager}
	general := domain.User{Role: domain.RoleGeneral}

	assert.True(t, resource.CanBeReservedBy(&manager))
	assert.False(t, resource.CanBeReservedBy(&general))
	assert.False(t, inactiveResource.CanBeReservedBy(&manager))
}
