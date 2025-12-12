package antiplagiat

type AntiPlagiarismEngine interface {
	Compare(a, b []byte) (float64, error)
}
