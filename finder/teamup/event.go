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
	Notes          interface{} `json:"notes"`
	Version        string      `json:"version"`
	Readonly       bool        `json:"readonly"`
	Tz             interface{} `json:"tz"`
	StartDt        string      `json:"start_dt"`
	EndDt          string      `json:"end_dt"`
	RistartDt      interface{} `json:"ristart_dt"`
	RsstartDt      interface{} `json:"rsstart_dt"`
	CreationDt     string      `json:"creation_dt"`
	UpdateDt       *string     `json:"update_dt"`
	DeleteDt       interface{} `json:"delete_dt"`
	Custom         Custom      `json:"custom"`
}
