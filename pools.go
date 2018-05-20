package log

import (
	"sync"
)

//Pools for things that we can be sure aren't handled by two things at once
type Pools struct {
	Nodes                    *sync.Pool
	PlainMessages            *sync.Pool
	ServiceMessages          *sync.Pool
	CompleteMessages         *sync.Pool
	GetResponses             *sync.Pool
	GetServiceResponses      *sync.Pool
	GetServiceLevelResponses *sync.Pool
}

var pools = &Pools{
	Nodes: &sync.Pool{
		New: func() interface{} {
			return &Node{}
		},
	},

	PlainMessages: &sync.Pool{
		New: func() interface{} {
			return &PlainMessage{}
		},
	},

	ServiceMessages: &sync.Pool{
		New: func() interface{} {
			return &ServiceMessage{}
		},
	},

	CompleteMessages: &sync.Pool{
		New: func() interface{} {
			return &CompleteMessage{}
		},
	},

	GetResponses: &sync.Pool{
		New: func() interface{} {
			return &GetResponse{}
		},
	},

	GetServiceResponses: &sync.Pool{
		New: func() interface{} {
			return &GetServiceResponse{}
		},
	},

	GetServiceLevelResponses: &sync.Pool{
		New: func() interface{} {
			return &GetServiceLevelResponse{}
		},
	},
}
