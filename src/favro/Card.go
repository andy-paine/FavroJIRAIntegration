package favro

type FavroResponse struct {
	Cards     []Card `json:"entities"`
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	Pages     int    `json:"pages"`
	RequestID string `json:"requestId"`
}

type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Card struct {
	CardID              string        `json:"cardId"`
	CardCommonID        string        `json:"cardCommonId"`
	OrganizationID      string        `json:"organizationId"`
	Archived            bool          `json:"archived"`
	Position            float32       `json:"position"`
	Name                string        `json:"name"`
	WidgetCommonID      string        `json:"widgetCommonId"`
	ColumnID            string        `json:"columnId"`
	IsLane              bool          `json:"isLane"`
	ParentCardID        interface{}   `json:"parentCardId"`
	DetailedDescription string        `json:"detailedDescription"`
	Tags                []interface{} `json:"tags"`
	SequentialID        int           `json:"sequentialId"`
	TasksTotal          int           `json:"tasksTotal"`
	TasksDone           int           `json:"tasksDone"`
}

type PostCard struct {
	WidgetCommonID      string `json:"widgetCommonId"`
	Name                string `json:"name"`
	DetailedDescription string `json:"detailedDescription"`
	Tags                []Tag  `json:"tags"`
}
