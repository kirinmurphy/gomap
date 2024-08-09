package locationManager

type LocationDiff struct {
	Added   []Location `json:"added"`
	Removed []Location `json:"removed"`
	Changed []Location `json:"changed"`
}

func DiffLocations(oldLocations, newLocations []Location) LocationDiff {
	oldMap := make(map[string]Location)
	newMap := make(map[string]Location)
	diff := LocationDiff{}

	for _, loc := range oldLocations {
		oldMap[loc.Name] = loc
	}

	for _, loc := range newLocations {
		newMap[loc.Name] = loc
	}

	for name, newLoc := range newMap {
		if oldLoc, exists := oldMap[name]; exists {
			if newLoc != oldLoc {
				diff.Changed = append(diff.Changed, newLoc)
			}
			delete(oldMap, name)
		} else {
			diff.Added = append(diff.Added, newLoc)
		}
	}

	for _, oldLoc := range oldMap {
		diff.Removed = append(diff.Removed, oldLoc)
	}

	return diff
}
