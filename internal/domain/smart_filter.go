package domain

const keepRecent = 10

// SmartFilter aplica as regras de negócio para gerar o índice otimizado.
type SmartFilter struct{}

func NewSmartFilter() *SmartFilter {
	return &SmartFilter{}
}

func (f *SmartFilter) Apply(collection ReleaseCollection) ReleaseCollection {
	if collection.IsEmpty() {
		return collection
	}

	limit := f.calculateLimit(len(collection.Releases))
	filtered := make([]Release, 0)
	seenMajors := make(map[string]bool)

	filtered = f.appendRecent(filtered, collection.Releases, limit, seenMajors)
	filtered = f.appendLegacy(filtered, collection.Releases, limit, seenMajors)

	return ReleaseCollection{Releases: filtered}
}

func (f *SmartFilter) calculateLimit(total int) int {
	if total < keepRecent {
		return total
	}
	return keepRecent
}

func (f *SmartFilter) appendRecent(filtered []Release, all []Release, limit int, seen map[string]bool) []Release {
	for i := 0; i < limit; i++ {
		filtered = append(filtered, all[i])
		f.markSeen(seen, all[i].Major)
	}
	return filtered
}

func (f *SmartFilter) appendLegacy(filtered []Release, all []Release, start int, seen map[string]bool) []Release {
	for i := start; i < len(all); i++ {
		filtered = f.processLegacyItem(filtered, all[i], seen)
	}
	return filtered
}

func (f *SmartFilter) processLegacyItem(filtered []Release, item Release, seen map[string]bool) []Release {
	if f.shouldKeepLegacy(item.Major, seen) {
		filtered = append(filtered, item)
		f.markSeen(seen, item.Major)
	}
	return filtered
}

func (f *SmartFilter) shouldKeepLegacy(major string, seen map[string]bool) bool {
	if major == "" {
		return false
	}
	if major == "Outros" {
		return false
	}
	return !seen[major]
}

func (f *SmartFilter) markSeen(seen map[string]bool, major string) {
	if major != "" {
		seen[major] = true
	}
}
