package models

// Permission represents a specific action that can be performed in the system
type Permission string

// Permission constants - Define all available permissions in the system
const (
	// Product Management Permissions
	PermissionViewProducts   Permission = "view:products"
	PermissionCreateProducts Permission = "create:products"
	PermissionUpdateProducts Permission = "update:products"
	PermissionDeleteProducts Permission = "delete:products"
	PermissionManageStock    Permission = "manage:stock"

	// Order Management Permissions
	PermissionViewOwnOrders Permission = "view:own_orders"
	PermissionViewAllOrders Permission = "view:all_orders"
	PermissionCreateOrders  Permission = "create:orders"
	PermissionUpdateOrders  Permission = "update:orders"
	PermissionCancelOrders  Permission = "cancel:orders"

	// User Management Permissions
	PermissionViewUsers   Permission = "view:users"
	PermissionCreateUsers Permission = "create:users"
	PermissionUpdateUsers Permission = "update:users"
	PermissionDeleteUsers Permission = "delete:users"
	PermissionManageRoles Permission = "manage:roles"

	// Media Management Permissions
	PermissionUploadMedia Permission = "upload:media"
	PermissionDeleteMedia Permission = "delete:media"

	// Notification Permissions
	PermissionSendNotifications Permission = "send:notifications"

	// Category Management Permissions
	PermissionManageCategories Permission = "manage:categories"

	// Variant Management Permissions
	PermissionManageVariants Permission = "manage:variants"

	// Pricing Permissions
	PermissionManagePricing Permission = "manage:pricing"
	PermissionViewPricing   Permission = "view:pricing"
)

// RolePermissions maps each role to its set of permissions
var RolePermissions = map[UserRole][]Permission{
	RoleAdmin: {
		// Product permissions
		PermissionViewProducts,
		PermissionCreateProducts,
		PermissionUpdateProducts,
		PermissionDeleteProducts,
		PermissionManageStock,

		// Order permissions
		PermissionViewOwnOrders,
		PermissionViewAllOrders,
		PermissionCreateOrders,
		PermissionUpdateOrders,
		PermissionCancelOrders,

		// User permissions
		PermissionViewUsers,
		PermissionCreateUsers,
		PermissionUpdateUsers,
		PermissionDeleteUsers,
		PermissionManageRoles,

		// Media permissions
		PermissionUploadMedia,
		PermissionDeleteMedia,

		// Notification permissions
		PermissionSendNotifications,

		// Category permissions
		PermissionManageCategories,

		// Variant permissions
		PermissionManageVariants,

		// Pricing permissions
		PermissionManagePricing,
		PermissionViewPricing,
	},
	RoleCustomer: {
		// Product permissions (read-only)
		PermissionViewProducts,

		// Order permissions (own orders only)
		PermissionViewOwnOrders,
		PermissionCreateOrders,

		// Pricing permissions (read-only)
		PermissionViewPricing,
	},
}

// HasPermission checks if a given role has a specific permission
func HasPermission(role UserRole, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetRolePermissions returns all permissions for a given role
func GetRolePermissions(role UserRole) []Permission {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []Permission{}
	}
	return permissions
}

// CanAccessResource checks if a user can access a specific resource
// This is used for ownership-based access control
func CanAccessResource(userID uint, resourceOwnerID uint, userRole UserRole) bool {
	// Admins can access all resources
	if userRole == RoleAdmin {
		return true
	}

	// Users can only access their own resources
	return userID == resourceOwnerID
}

// IsAdmin checks if a role is admin
func IsAdmin(role UserRole) bool {
	return role == RoleAdmin
}

// IsCustomer checks if a role is customer
func IsCustomer(role UserRole) bool {
	return role == RoleCustomer
}

// ValidateRole checks if a role is valid
func ValidateRole(role UserRole) bool {
	return role == RoleAdmin || role == RoleCustomer
}
