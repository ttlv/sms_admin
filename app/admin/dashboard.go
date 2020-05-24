package admin

import (
	"github.com/qor/admin"
)

// SetupDashboard setup dashboard
func SetupDashboard(Admin *admin.Admin) {
	// Add Dashboard
	Admin.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin", Priority: 1})
}
