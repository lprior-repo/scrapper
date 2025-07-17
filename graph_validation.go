package main

import (
	"fmt"
	"strconv"
)

// validateNodeCreation validates node creation parameters (Pure Core)
func validateNodeCreation(request GraphOperationRequest) error {
	if request.Label == "" {
		return fmt.Errorf("node label is required")
	}
	if request.Properties == nil {
		return fmt.Errorf("node properties cannot be nil")
	}
	return nil
}

// validateNodeRetrieval validates node retrieval parameters (Pure Core)
func validateNodeRetrieval(request GraphOperationRequest) error {
	if request.NodeID == "" {
		return fmt.Errorf("node ID is required")
	}
	if !validateNodeID(request.NodeID) {
		return fmt.Errorf("invalid node ID format")
	}
	return nil
}

// validateRelationshipCreation validates relationship creation parameters (Pure Core)
func validateRelationshipCreation(request GraphOperationRequest) error {
	if request.FromID == "" {
		return fmt.Errorf("from node ID is required")
	}
	if request.ToID == "" {
		return fmt.Errorf("to node ID is required")
	}
	if request.RelType == "" {
		return fmt.Errorf("relationship type is required")
	}
	if !validateNodeID(request.FromID) {
		return fmt.Errorf("invalid from node ID format")
	}
	if !validateNodeID(request.ToID) {
		return fmt.Errorf("invalid to node ID format")
	}
	return nil
}

// validateRelationshipUpdate validates relationship update parameters (Pure Core)
func validateRelationshipUpdate(request GraphOperationRequest) error {
	if request.NodeID == "" {
		return fmt.Errorf("relationship ID is required")
	}
	if !validateNodeID(request.NodeID) {
		return fmt.Errorf("invalid relationship ID format")
	}
	if request.Properties == nil {
		return fmt.Errorf("properties cannot be nil")
	}
	return nil
}

// validateRelationshipType validates relationship type against business rules (Pure Core)
func validateRelationshipType(relType string) error {
	if relType == "" {
		return fmt.Errorf("relationship type cannot be empty")
	}

	// Define allowed relationship types
	allowedTypes := map[string]bool{
		"KNOWS":          true,
		"WORKS_WITH":     true,
		"MANAGES":        true,
		"REPORTS_TO":     true,
		"DEPENDS_ON":     true,
		"CONNECTED_TO":   true,
		"CONTAINS":       true,
		"BELONGS_TO":     true,
		"REFERENCES":     true,
		"FOLLOWS":        true,
		"LIKES":          true,
		"OWNS":           true,
		"RELATES_TO":     true,
		"CREATED":        true,
		"MODIFIED":       true,
		"HAS_PERMISSION": true,
		"IS_MEMBER_OF":   true,
		"APPROVED_BY":    true,
		"REVIEWED_BY":    true,
		"HAS_CODEOWNER":  true,
		"MAINTAINS":      true,
		"CONTRIBUTES_TO": true,
	}

	if !allowedTypes[relType] {
		return fmt.Errorf("invalid relationship type: %s", relType)
	}

	return nil
}

// validateRelationshipProperties validates relationship properties against business rules (Pure Core)
func validateRelationshipProperties(relType string, properties map[string]interface{}) error {
	if properties == nil {
		return nil // Properties are optional
	}

	// Common property validations
	if strength, exists := properties["strength"]; exists {
		if strengthVal, ok := strength.(float64); ok {
			if strengthVal < 0.0 || strengthVal > 1.0 {
				return fmt.Errorf("strength must be between 0.0 and 1.0")
			}
		}
	}

	if weight, exists := properties["weight"]; exists {
		if weightVal, ok := weight.(float64); ok {
			if weightVal < 0.0 {
				return fmt.Errorf("weight must be non-negative")
			}
		}
	}

	// Type-specific validations
	switch relType {
	case "KNOWS":
		if since, exists := properties["since"]; exists {
			if sinceStr, ok := since.(string); ok {
				if len(sinceStr) < 4 { // At least year format
					return fmt.Errorf("since date must be at least 4 characters")
				}
			}
		}

	case "MANAGES", "REPORTS_TO":
		if level, exists := properties["level"]; exists {
			if levelVal, ok := level.(float64); ok {
				if levelVal < 1.0 || levelVal > 10.0 {
					return fmt.Errorf("management level must be between 1 and 10")
				}
			}
		}

	case "HAS_PERMISSION":
		if permission, exists := properties["permission"]; exists {
			if permStr, ok := permission.(string); ok {
				allowedPermissions := map[string]bool{
					"read": true, "write": true, "delete": true, "admin": true,
				}
				if !allowedPermissions[permStr] {
					return fmt.Errorf("invalid permission type: %s", permStr)
				}
			}
		}
	}

	return nil
}

// validateRelationshipBusinessRules validates relationship against business rules (Pure Core)
func validateRelationshipBusinessRules(request GraphOperationRequest) error {
	// Validate relationship type
	if err := validateRelationshipType(request.RelType); err != nil {
		return err
	}

	// Validate properties
	if err := validateRelationshipProperties(request.RelType, request.Properties); err != nil {
		return err
	}

	// Prevent self-relationships for certain types
	if request.FromID == request.ToID {
		restrictedTypes := map[string]bool{
			"MANAGES":    true,
			"REPORTS_TO": true,
		}
		if restrictedTypes[request.RelType] {
			return fmt.Errorf("self-relationships not allowed for type: %s", request.RelType)
		}
	}

	return nil
}

// validateNodeID checks if a node ID is valid (Pure Core)
func validateNodeID(id string) bool {
	if id == "" {
		return false
	}
	_, err := strconv.ParseInt(id, 10, 64)
	return err == nil
}
