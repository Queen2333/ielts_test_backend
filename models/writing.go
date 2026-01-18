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
	Status    			int    				`json:"status"`
	Type      			int    				`json:"type"`
	PartList			[]int				`json:"part_list"`
	UserID				string				`json:"user_id,omitempty"`
}

type WritingPartItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Type				string				`json:"type,omitempty"`
	Source				string				`json:"source,omitempty"`
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