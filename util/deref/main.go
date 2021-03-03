package deref

func String(s *string, or string) string {
	if s == nil {
		return or
	}
	return *s
}
