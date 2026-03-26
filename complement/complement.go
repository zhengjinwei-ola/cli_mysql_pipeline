package complement

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/model"
	"time"
)

var Instance *Forward

func init() {
	Instance = &Forward{
		data: make([]*model.OriginRow, 0),
	}
	go Instance.Start()
}

type Forward struct {
	data     []*model.OriginRow
	Receiver chan *model.OriginRow
}

func (f *Forward) Add(message *model.OriginRow) {
	f.data = append(f.data, message)
}

func (f *Forward) Start() {
	timer := time.NewTicker(time.Millisecond * 100)
	for range timer.C {
		if len(f.data) > 0 {
			for i := 0; i < len(f.data); i++ {
				f.Receiver <- f.data[i]
			}
			f.data = make([]*model.OriginRow, 0)
		}
	}
}
