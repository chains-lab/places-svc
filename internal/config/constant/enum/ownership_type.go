package enum

import "fmt"

const (
	OwnershipTypePrivate   = "private"   // Private ownership type
	OwnershipTypePublic    = "public"    // Public ownership type
	OwnershipTypeMunicipal = "municipal" // Municipal ownership type
)

var ownershipTypes = []string{
	OwnershipTypePrivate,
	OwnershipTypePublic,
	OwnershipTypeMunicipal,
}

var ErrorOwnershipTypeNotSupported = fmt.Errorf("ownership type must be one of: %s", GetAllOwnershipTypes())

// ParseOwnershipType checks if the provided ownership type is valid and returns it.
func ParseOwnershipType(ownershipType string) (string, error) {
	for _, t := range ownershipTypes {
		if t == ownershipType {
			return t, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", ownershipType, ErrorOwnershipTypeNotSupported)
}

func GetAllOwnershipTypes() []string {
	return ownershipTypes
}
