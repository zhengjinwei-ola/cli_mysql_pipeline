package common

import (
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNsqMsgToMap(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		var s = []byte("a:39:{i:105699063;i:2;i:105699069;i:2;i:105699070;i:2;i:105699073;i:2;i:105699076;i:2;i:105699078;i:2;i:105699081;i:2;i:105699082;i:2;i:105699084;i:2;i:105699086;i:2;i:105699087;i:2;i:105699088;i:2;i:105699091;i:2;i:105699092;i:2;i:105699093;i:1;i:105699098;i:1;i:105699103;i:2;i:105699137;i:2;i:105699140;i:2;i:105699143;i:2;i:105699149;i:2;i:105699154;i:2;i:105699168;i:2;i:105699209;i:2;i:105699212;i:2;i:105699217;i:2;i:105699227;i:2;i:105699294;i:2;i:105699323;i:2;i:105699614;i:2;i:105699664;i:2;i:105699777;i:2;i:105700174;i:2;i:105700278;i:2;i:105700880;i:2;i:105702237;i:2;i:105703213;i:2;i:105706083;i:2;i:105707266;i:2;}")
		msg := &nsq.Message{Body: s}
		res, err := NsqMsgToMap(msg)
		assert.Nil(t, err)
		assert.IsType(t, map[string]interface{}{}, res)
	})
}
