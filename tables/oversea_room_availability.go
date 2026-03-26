package tables

var DefaultOverseaRoomAvailability map[string]interface{} = map[string]interface{}{
	"room_available_rid":  0,
	"room_available_seat": 0,
}

const (
	sqlChatroomConfigAvailability string = "select rid, SUM(CASE when uid > 0 and `lock` = 0 and forbidden = 0 then 1 else 0 end) as 'available' from xs_chatroom_config where rid = ? group by rid"
)

type OpRowOverseaRoomAvailability struct {
	OpRowTable
}
