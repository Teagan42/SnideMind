package extracttags

import (
	"testing"
)

func TestTagTree_UniqueIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, tag := range TagTree {
		if tag.ID == "" {
			t.Errorf("Tag with empty ID found: %+v", tag)
		}
		if seen[tag.ID] {
			t.Errorf("Duplicate tag ID found: %s", tag.ID)
		}
		seen[tag.ID] = true
	}
}

func TestTagTree_ParentIDsExist(t *testing.T) {
	idSet := make(map[string]struct{})
	for _, tag := range TagTree {
		idSet[tag.ID] = struct{}{}
	}
	for _, tag := range TagTree {
		if tag.ParentID != nil {
			if _, ok := idSet[*tag.ParentID]; !ok {
				t.Errorf("Tag %s has non-existent ParentID: %s", tag.ID, *tag.ParentID)
			}
		}
	}
}

func TestTagTree_DescriptionsAndNames(t *testing.T) {
	for _, tag := range TagTree {
		if tag.Name == "" {
			t.Errorf("Tag %s has empty Name", tag.ID)
		}
		if tag.Description == "" {
			t.Errorf("Tag %s has empty Description", tag.ID)
		}
	}
}

func TestTagTree_ParentChildRelationship(t *testing.T) {
	// For each tag with a ParentID, ensure the parent does not have a ParentID that is the same as the child
	for _, tag := range TagTree {
		if tag.ParentID != nil {
			for _, parent := range TagTree {
				if parent.ID == *tag.ParentID && parent.ParentID != nil && *parent.ParentID == tag.ID {
					t.Errorf("Circular parent-child relationship between %s and %s", tag.ID, parent.ID)
				}
			}
		}
	}
}
