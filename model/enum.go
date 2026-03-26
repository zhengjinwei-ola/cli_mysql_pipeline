package model

type OpType int8
type FlowType int8
type EsType int8

const (
	OpType_Write      OpType = 0
	OpType_Update     OpType = 1
	OpType_Delete     OpType = 2
	OpType_Update_Row OpType = 3
)

const (
	FlowType_Main   FlowType = 0 //主
	FlowType_Join   FlowType = 1 //根据主键展开合并
	FlowType_Array  FlowType = 2 //nested 类型
	FlowType_Append FlowType = 3 //直接附加到某个字段上
)

const (
	EsType_Number      EsType = 0
	EsType_Text        EsType = 1
	EsType_TureOrFalse EsType = 2
	EsType_SetNumber   EsType = 3
	EsType_SetString   EsType = 4
	EsType_Float       EsType = 5
	EsType_Func        EsType = 6
)

const (
	True_Upper  string = "TRUE"
	True_Lower  string = "true"
	False_Upper string = "FALSE"
	False_Lower string = "false"
)
