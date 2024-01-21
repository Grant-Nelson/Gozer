package tools

type Tool interface {
	Run(args []string) error
}
