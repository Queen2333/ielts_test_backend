package models

type WritingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status,omitempty"`
	Type				string				`json:"type,omitempty"`
	PartList			[]WritingPartItem	`json:"part_list"`
}

type BasicWritingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status    			FlexInt    			`json:"status"`
	Type      			FlexInt    			`json:"type"`
	PartList			[]int				`json:"part_list"`
	UserID				string				`json:"user_id,omitempty"`
}

type WritingPartItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Type				string				`json:"type,omitempty"`        // 数据来源：1=系统，2=官方，3=用户（与其他表统一）
	TaskType			string				`json:"task_type,omitempty"`   // 任务类型：1=Task1，2=Task2
	Title				string				`json:"title"`
	SubTitle			string				`json:"sub_title,omitempty"`
	Img					string				`json:"img,omitempty"`
	UserID				string				`json:"user_id,omitempty"`
}
type WritingRecordsItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status"`
	Type				string				`json:"type"`
	Answers				[]interface{}		`json:"answers"`
	UserID				string				`json:"user_id,omitempty"`
	RestSeconds			int					`json:"rest_seconds,omitempty"`
}