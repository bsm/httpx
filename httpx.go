package httpx

func coalesce(vv ...string) string {
	for _, v := range vv {
		if v != "" {
			return v
		}
	}
	return ""
}
