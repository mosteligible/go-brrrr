package loads

import (
	"sync"

	"github.com/mosteligible/go-brrrr/pkg/types"
)

type Loader interface {
	Load(
		parameters *types.Parameters,
		mq *types.MetricsQueue,
		wg *sync.WaitGroup,
		timeStr string,
	) error
}
