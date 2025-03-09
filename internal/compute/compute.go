package compute

type Compute interface {
	Parse(data string) (int, []string, error)
}
