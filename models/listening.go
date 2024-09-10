package models

type ListeningItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status				string				`json:"status,omitempty"`
	AudioFiles			[]string			`json:"audio_files"`
	Type				string				`json:"type,omitempty"`
	PartList			[]ListeningPartItem	`json:"part_list"`
}

type BasicListeningItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	Status    			int    				`json:"status"`
	Type      			int    				`json:"type"`
	AudioFiles 			[]string			`json:"audio_files"`
	PartList			[]int				`json:"part_list"`
	UserID				string				`json:"user_id,omitempty"`
}

type ListeningPartItem struct {
	ID 					int					`json:"id,omitempty"`
	Name				string				`json:"name"`
	TypeList			[]ListeningTypeItem	`json:"type_list"`
	Type				string				`json:"type,omitempty"`
}

type ListeningTypeItem struct {
	Title    			string    				`json:"title"`
	Type     			string       			`json:"type"`
	ArticleContent    	string					`json:"article_content,omitempty"`
	Picture				[]PicturesItem			`json:"picture,omitempty"`
	MatchingOptions     []MatchingOptionsItem	`json:"matching_options,omitempty"`
	QuestionList		[]ListeningQuestionItem	`json:"question_list"`
}

type ListeningQuestionItem struct {
	No                  string			`json:"no"`
	Question			string			`json:"question"`
	Options				[]OptionsItem	`json:"options,omitempty"`
	Answer				interface{}		`json:"answer"`
}

type OptionsItem struct {
	Label				string			`json:"label"`
	Value				string			`json:"value"`
}

type MatchingOptionsItem struct {
	Label				string			`json:"label"`
	Content				string			`json:"content"`
}

type PicturesItem struct {
	Url				string			`json:"url"`
	Name			string			`json:"name,omitempty"`
}