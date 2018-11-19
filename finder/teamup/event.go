package teamup

type Event struct {
	ID             string      `json:"id"`
	SeriesID       interface{} `json:"series_id"`
	RemoteID       interface{} `json:"remote_id"`
	SubcalendarID  int64       `json:"subcalendar_id"`
	SubcalendarIDS []int64     `json:"subcalendar_ids"`
	AllDay         bool        `json:"all_day"`
	Rrule          string      `json:"rrule"`
	Title          string      `json:"title"`
	Who            string      `json:"who"`
	Location       string      `json:"location"`
	Notes          string      `json:"notes"`
	Version        string      `json:"version"`
	Readonly       bool        `json:"readonly"`
	Tz             interface{} `json:"tz"`
	StartDt        DateTime    `json:"start_dt"`
	EndDt          DateTime    `json:"end_dt"`
	RistartDt      interface{} `json:"ristart_dt"`
	RsstartDt      interface{} `json:"rsstart_dt"`
	CreationDt     DateTime    `json:"creation_dt"`
	UpdateDt       *DateTime   `json:"update_dt"`
	DeleteDt       interface{} `json:"delete_dt"`
}
