package bigarea

import (
	"github.com/olachat/banban_server/cli_mysql_pipeline/tools/batch_importer/bigarea/mocks"
	"github.com/olachat/banban_server/common/serialize"
	"errors"
	"github.com/golang/mock/gomock"
	goNsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_wrapBigAreaHandlers(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  string
	}{
		{
			name: "happy case",
			input: map[string]string{
				"type": "room.bigarea.batch",
				"data": "something",
			},
			want: "room",
		},
		{
			name: "happy case",
			input: map[string]string{
				"type": "user.bigarea.batch",
				"data": "something",
			},
			want: "user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			buf, _ := serialize.Marshal(tt.input)
			msg := &goNsq.Message{Body: buf}

			buf, _ = serialize.Marshal(tt.input["data"])
			hRoom := mocks.NewMockiBatchImportHandler(ctrl)
			hRoom.EXPECT().HandleMsg(&goNsq.Message{Body: buf}).
				Return(errors.New("room" + string(buf))).AnyTimes()
			hUser := mocks.NewMockiBatchImportHandler(ctrl)
			hUser.EXPECT().HandleMsg(&goNsq.Message{Body: buf}).
				Return(errors.New("user" + string(buf))).AnyTimes()

			f := WrapBigAreaHandlers(hRoom, hUser)
			err := f(msg)
			s := err.Error()
			assert.Equal(t, tt.want+string(buf), s)
		})
	}
}

func Test_nsqMsg_getData(t *testing.T) {
	t.Run("testing unmarshal", func(t *testing.T) {
		var s = []byte("a:2:{s:4:\"type\";s:18:\"room.bigarea.batch\";s:4:\"data\";a:29:{i:105700351;i:2;i:105700352;i:2;i:105700360;i:2;i:105700364;i:2;i:105700368;i:2;i:105700369;i:2;i:105700374;i:2;i:105700378;i:2;i:105700381;i:2;i:105700385;i:2;i:105700390;i:2;i:105700394;i:2;i:105700403;i:2;i:105700448;i:2;i:105700551;i:2;i:105700587;i:2;i:105700692;i:2;i:105700810;i:2;i:105701025;i:2;i:105701412;i:2;i:105701549;i:2;i:105701911;i:2;i:105702292;i:2;i:105702493;i:2;i:105703250;i:2;i:105703530;i:2;i:105704293;i:2;i:105705338;i:2;i:105706734;i:2;}}")
		body, err := msgBody(&goNsq.Message{Body: s})
		if err != nil {
			panic(err)
		}
		msg := body.getData()
		println(msg)
	})
}

func Test_msgBody(t *testing.T) {
	t.Run("testing wrong type", func(t *testing.T) {
		var s = []byte("a:2:{s:4:\"type\";s:18:\"bigarea.batch\";s:4:\"data\";a:29:{i:105700351;i:2;i:105700352;i:2;i:105700360;i:2;i:105700364;i:2;i:105700368;i:2;i:105700369;i:2;i:105700374;i:2;i:105700378;i:2;i:105700381;i:2;i:105700385;i:2;i:105700390;i:2;i:105700394;i:2;i:105700403;i:2;i:105700448;i:2;i:105700551;i:2;i:105700587;i:2;i:105700692;i:2;i:105700810;i:2;i:105701025;i:2;i:105701412;i:2;i:105701549;i:2;i:105701911;i:2;i:105702292;i:2;i:105702493;i:2;i:105703250;i:2;i:105703530;i:2;i:105704293;i:2;i:105705338;i:2;i:105706734;i:2;}}")
		_, err := msgBody(&goNsq.Message{Body: s})
		assert.NotNil(t, err)
	})
	t.Run("testing wrong type", func(t *testing.T) {
		var s = []byte("a:2:{s:4:\"type\";s:18:\"room.bigarea.batch\";s:4:\"data\";a:29:{i:105700351;i:2;i:105700352;i:2;i:105700360;i:2;i:105700364;i:2;i:105700368;i:2;i:105700369;i:2;i:105700374;i:2;i:105700378;i:2;i:105700381;i:2;i:105700385;i:2;i:105700390;i:2;i:105700394;i:2;i:105700403;i:2;i:105700448;i:2;i:105700551;i:2;i:105700587;i:2;i:105700692;i:2;i:105700810;i:2;i:105701025;i:2;i:105701412;i:2;i:105701549;i:2;i:105701911;i:2;i:105702292;i:2;i:105702493;i:2;i:105703250;i:2;i:105703530;i:2;i:105704293;i:2;i:105705338;i:2;i:105706734;i:2;}}")
		_, err := msgBody(&goNsq.Message{Body: s})
		assert.Nil(t, err)
	})
}
