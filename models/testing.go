package models

type TestingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status,omitempty"`
	ReadingParts		[]ReadingPartItem	`json:"reading_parts"`
	ListeningParts		[]ListeningPartItem `json:"listening_parts"`
	WritingParts		[]WritingPartItem	`json:"writing_parts"`
	Type				string				`json:"type"`
}

type BasicTestingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status    			FlexInt    			`json:"status"`
	Type      			FlexInt    			`json:"type"`
	ListeningIDs		[]int				`json:"listening_ids"`
	ReadingIDs			[]int				`json:"reading_ids"`
	WritingIDs			[]int				`json:"writing_ids"`
	UserID				string				`json:"user_id,omitempty"`
}
type TestingRecordsItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status"`
	Type				string				`json:"type"`
	Score				[]int				`json:"score,omitempty"`
	Answers				[]AnswerItem		`json:"answers"`
	UserID				string				`json:"user_id,omitempty"`
	RestSeconds			[]int				`json:"rest_seconds,omitempty"`
	TestID 				int					`json:"test_id"`
}

type TestingAnswerItem struct {
	PartName			string				`json:"part_name"`
	AnswerList			[]interface{}		`json:"answer_list"`
}