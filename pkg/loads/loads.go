package loads

import (
	"go-brrrr/pkg/types"
	"sync"
)

type Loader interface {
	Load(
		parameters map[string]string,
		mq *types.MetricsQueue,
		wg *sync.WaitGroup,
	) error
}
