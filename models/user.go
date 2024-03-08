package models

type UserAuthPayload struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	StudentId string `json:"studentId"`
}

type UserCookies struct {
	Antiforgery string `json:"antiforgery"`
	Session     string `json:"session"`
	Identity    string `json:"identity"`
}

type UserResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type PersonalUserData struct {
	Username    string `json:"username"`
	Surname     string `json:"surname"`
	OtherNames  string `json:"other_names"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	Religion    string `json:"religion"`
}

type ProgrammeUserData struct {
	StudentId       string `json:"studentId"`
	IndexNo         string `json:"indexNo"`
	ProgrammeStream string `json:"programme_stream"`
}

type ContactUserData struct {
	SchoolEmail        string `json:"school_email"`
	PersonalEmail      string `json:"personal_email"`
	KNUSTMobile        string `json:"knust_mobile"`
	PersonalMobile     string `json:"personal_mobile"`
	AltPersonalMobile  string `json:"alt_personal_mobile"`
	PostalAddress      string `json:"postal_address"`
	ResidentialAddress string `json:"residential_address"`
}

type UserData struct {
	Personal  PersonalUserData  `json:"personal"`
	Programme ProgrammeUserData `json:"programme"`
	Contact   ContactUserData   `json:"contact"`
}

type UserDataResponse struct {
	Message  string   `json:"message"`
	UserData UserData `json:"user_data"`
}

type ResultsSelection struct {
	Years []string `json:"years"`
	Sems  []string `json:"sems"`
}

type ResultsSelectionResponse struct {
	Message string           `json:"message"`
	Results ResultsSelection `json:"results"`
}

type GetResultsPayload struct {
	Year string `json:"year"`
	Sem  string `json:"sem"`
}

type ResultsPersonalData struct {
	Name      string `json:"name"`
	Year      string `json:"year"`
	IndexNo   string `json:"indexNo"`
	Programme string `json:"programme"`
	StudentID string `json:"studentId"`
	Date      string `json:"date"`
	Option    string `json:"option"`
	Username  string `json:"username"`
}

type Results struct {
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Credits    string `json:"credits"`
	Grade      string `json:"grade"`
	TotalMark  string `json:"total_mark"`
}

type ResultsSummaryExtra struct {
	Semester   string `json:"semester"`
	Cumulative string `json:"cumulative"`
}

type ResultsSummary struct {
	CreditsRegistered ResultsSummaryExtra `json:"credits_registered"`
	CreditsObtained   ResultsSummaryExtra `json:"credits_obtained"`
	CreditsCalculated ResultsSummaryExtra `json:"credits_calculated"`
	WeightedMarks     ResultsSummaryExtra `json:"weighted_marks"`
	CWA               ResultsSummaryExtra `json:"cwa"`
}

type GetResultsResponse struct {
	Message      string              `json:"message"`
	PersonalData ResultsPersonalData `json:"personal_data"`
	Results      []Results           `json:"results"`
	Summary      ResultsSummary      `json:"summary"`
	Trails       []string            `json:"trails"`
}
