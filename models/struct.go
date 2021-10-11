package models

type APIResponse struct {
	Professors         []Professors `json:"professors"`
	SearchResultsTotal int          `json:"searchResultsTotal"`
	Remaining          int          `json:"remaining"`
	Type               string       `json:"type"`
}
type Professors struct {
	TDept           string `json:"tDept"`
	TSid            string `json:"tSid"`
	InstitutionName string `json:"institution_name"`
	TFname          string `json:"tFname"`
	TMiddlename     string `json:"tMiddlename"`
	TLname          string `json:"tLname"`
	Tid             int    `json:"tid"`
	TNumRatings     int    `json:"tNumRatings"`
	RatingClass     string `json:"rating_class"`
	ContentType     string `json:"contentType"`
	CategoryType    string `json:"categoryType"`
	OverallRating   string `json:"overall_rating"`
}

func (p *Professors) SetProfessor(professor Professors) {
	p.TDept = professor.TDept
	p.TSid = professor.TSid
	p.InstitutionName = professor.InstitutionName
	p.TFname = professor.TFname
	p.TMiddlename = professor.TMiddlename
	p.TLname = professor.TLname
	p.TMiddlename = professor.TMiddlename
	p.TFname = professor.TFname
	p.Tid = professor.Tid
	p.TNumRatings = professor.TNumRatings
	p.RatingClass = professor.RatingClass
	p.ContentType = professor.ContentType
	p.CategoryType = professor.CategoryType
	p.OverallRating = professor.OverallRating
}

type ProfessorRating struct {
	Ratings   []Ratings `json:"ratings"`
	Remaining int       `json:"remaining"`
}
type Ratings struct {
	Attendance        string      `json:"attendance"`
	ClarityColor      string      `json:"clarityColor"`
	EasyColor         string      `json:"easyColor"`
	HelpColor         string      `json:"helpColor"`
	HelpCount         int         `json:"helpCount"`
	ID                int         `json:"id"`
	NotHelpCount      int         `json:"notHelpCount"`
	OnlineClass       string      `json:"onlineClass"`
	Quality           string      `json:"quality"`
	RClarity          int         `json:"rClarity"`
	RClass            string      `json:"rClass"`
	RComments         string      `json:"rComments"`
	RDate             string      `json:"rDate"`
	REasy             float64     `json:"rEasy"`
	REasyString       string      `json:"rEasyString"`
	RErrorMsg         interface{} `json:"rErrorMsg"`
	RHelpful          int         `json:"rHelpful"`
	RInterest         string      `json:"rInterest"`
	ROverall          float64     `json:"rOverall"`
	ROverallString    string      `json:"rOverallString"`
	RStatus           int         `json:"rStatus"`
	RTextBookUse      string      `json:"rTextBookUse"`
	RTimestamp        int64       `json:"rTimestamp"`
	RWouldTakeAgain   string      `json:"rWouldTakeAgain"`
	SID               int         `json:"sId"`
	TakenForCredit    string      `json:"takenForCredit"`
	Teacher           interface{} `json:"teacher"`
	TeacherGrade      string      `json:"teacherGrade"`
	TeacherRatingTags []string    `json:"teacherRatingTags"`
	UnUsefulGrouping  string      `json:"unUsefulGrouping"`
	UsefulGrouping    string      `json:"usefulGrouping"`
}
