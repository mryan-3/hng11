package controller
import (
	"testing"

	"github.com/google/uuid"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/stretchr/testify/assert"
)


func TestOrganizationAccess(t *testing.T) {
	// Create test users and organizations
	user1 := &models.User{UserID: uuid.New()}
	user2 := &models.User{UserID: uuid.New()}
	org1 := &models.Organisation{ID: uuid.New(), Name: "Org 1"}
	org2 := &models.Organisation{ID: uuid.New(), Name: "Org 2"}

	// Associate user1 with org1, and user2 with org2
	user1.Organisations = []*models.Organisation{org1}
	user2.Organisations = []*models.Organisation{org2}

	// Test user1's access
	assert.True(t, hasAccessToOrg(user1, org1.ID))
	assert.False(t, hasAccessToOrg(user1, org2.ID))

	// Test user2's access
	assert.False(t, hasAccessToOrg(user2, org1.ID))
	assert.True(t, hasAccessToOrg(user2, org2.ID))
}

// Helper function to check if a user has access to an organization
func hasAccessToOrg(user *models.User, orgID uuid.UUID) bool {
	for _, org := range user.Organisations {
		if org.ID == orgID {
			return true
		}
	}
	return false
}
