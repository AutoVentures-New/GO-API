package pkg

type UlidGetter interface {
	GetUlid() string
}

func ExtractIdentifiers[T UlidGetter](items []T) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.GetUlid())
	}
	return ids
}

func ExtractField[T any, R any](items []T, extractor func(T) R) []R {
	results := make([]R, 0, len(items))
	for _, item := range items {
		results = append(results, extractor(item))
	}
	return results
}
