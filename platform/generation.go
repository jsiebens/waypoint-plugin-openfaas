package platform

import "context"

// GenerationFunc implements component.Generation
func (p *Platform) GenerationFunc() interface{} {
	return p.generation
}

// Generation returns the generation ID.
func (p *Platform) generation(ctx context.Context, ) ([]byte, error) {
	// Static generation since we will always use the `function_name` to
	// automatically delete the function.
	return []byte("openfaas"), nil
}
