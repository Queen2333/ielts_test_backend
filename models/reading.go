package models

type ReadingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status,omitempty"`
	Type				string				`json:"type,omitempty"`
	PartList			[]ReadingPartItem	`json:"part_list"`
}

type BasicReadingItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status    			FlexInt    			`json:"status"`
	Type      			FlexInt    			`json:"type"`
	PartList			[]int				`json:"part_list"`
	UserID				string				`json:"user_id,omitempty"`
}

type ReadingPartItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	TypeList			[]ReadingTypeItem	`json:"type_list"`
	Type				string				`json:"type,omitempty"`
	Article				string				`json:"article"`
	UserID				string				`json:"user_id,omitempty"`
}

type ReadingTypeItem struct {
	Title    			string    					`json:"title,omitempty"`
	Type     			string       				`json:"type"`
	NB					bool						`json:"nb,omitempty"`
	QuestionList		[]ReadingQuestionItem		`json:"question_list"`
	Options				[]ReadingOptionsItem		`json:"options,omitempty"`
	ArticleContent		string						`json:"article_content,omitempty"`
	Picture 			[]PicturesItem				`json:"picture,omitempty"`
	MatchingOptions     []MatchingOptionsItem		`json:"matching_options,omitempty"`
}

type ReadingQuestionItem struct {
	ID					int							`json:"id"`
	No                  string						`json:"no"`
	Question			string						`json:"question"`
	Options				[]QuestionOptionsItem		`json:"options,omitempty"`
	Answer				interface{}					`json:"answer"`
	AnswerCount			int							`json:"answer_count,omitempty"`
	Content				string						`json:"content,omitempty"`
	MatchedOption		*MatchingOptionsItem		`json:"matchedOption,omitempty"`
	IsDraggingOver		bool						`json:"isDraggingOver,omitempty"`
}

type ReadingOptionsItem struct {
	Label				string			`json:"label"`
	Content				string			`json:"content,omitempty"`
	ID					string			`json:"id"`
}

type QuestionOptionsItem struct {
	Label				string			`json:"label"`
	Text				string			`json:"text"`
}

type ReadingRecordsItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status"`
	Type				string				`json:"type"`
	Score				int					`json:"score,omitempty"`
	Answers				[]AnswerItem		`json:"answers"`
	UserID				string				`json:"user_id,omitempty"`
	RestSeconds			int					`json:"rest_seconds,omitempty"`
	TestID 				int					`json:"test_id"`
}
