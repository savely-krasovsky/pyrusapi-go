package pyrus

import (
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// TaskRequest is necessary to create a task.
type TaskRequest struct {
	Text                 string        `json:"text,omitempty"`
	Responsible          *Person       `json:"responsible,omitempty"`
	DueDate              string        `json:"due_date,omitempty"`
	Due                  *time.Time    `json:"due,omitempty"`
	Duration             int           `json:"duration,omitempty"`
	Subject              string        `json:"subject,omitempty"`
	Participants         []*Person     `json:"participants,omitempty"`
	Subscribers          []*Person     `json:"subscribers,omitempty"`
	ParentTaskID         int           `json:"parent_task_id,omitempty"`
	ListIDs              []int         `json:"list_ids,omitempty"`
	Attachments          []*Attachment `json:"attachments,omitempty"`
	ScheduledDate        string        `json:"scheduled_date,omitempty"`
	ScheduledDatetimeUTC *time.Time    `json:"scheduled_datetime_utc,omitempty"`
	Approvals            [][]*Person   `json:"approvals,omitempty"`
	FormID               int           `json:"form_id,omitempty"`
	FillDefaults         bool          `json:"fill_defaults,omitempty"`
}

// Validate allows to validate request before sending.
func (r TaskRequest) Validate() error {
	const bothMsg = "use text or form_id, not both"
	const eitherMsg = "use either text or form_id"

	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Text, validation.
			When(r.FormID != 0, validation.Empty.Error(bothMsg)).
			Else(validation.Required.Error(eitherMsg))),
		validation.Field(&r.FormID, validation.
			When(r.Text != "", validation.Empty.Error(bothMsg)).
			Else(validation.Required.Error(eitherMsg))),
		validation.Field(&r.Due,
			validation.When(r.DueDate != "", validation.Nil.Error("use due or due_date, not both"), validation.Date("2006-01-02")),
			validation.When(r.Duration != 0, validation.Required.Error("duration requires due")),
		),
		validation.Field(&r.DueDate, validation.When(r.Due != nil, validation.Empty.Error("use due or due_date, not both"))),
		validation.Field(&r.Duration, validation.Min(0), validation.Max(365*24*60)),
		validation.Field(&r.Responsible),
		validation.Field(&r.Participants, validation.Each()),
		validation.Field(&r.Subscribers, validation.Each()),
		validation.Field(&r.Approvals, validation.Each()),
	)
}

// Validate allows to validate request before sending.
func (p Person) Validate() error {
	const bothMsg = "use id or email, not both"
	const eitherMsg = "use either id or email"

	return validation.ValidateStruct(
		&p,
		validation.Field(&p.ID, validation.When(p.Email != "", validation.Empty.Error(bothMsg)).Else(validation.Required.Error(eitherMsg))),
		validation.Field(&p.Email, validation.When(p.ID != 0, validation.Empty.Error(bothMsg)).Else(validation.Required.Error(eitherMsg)), is.Email),
	)
}

// Validate allows to validate request before sending.
func (f FormField) Validate() error {
	const bothMsg = "use id or name, not both"
	const eitherMsg = "use either id or name"

	return validation.ValidateStruct(
		&f,
		validation.Field(&f.ID, validation.When(f.Name != "", validation.Empty.Error(bothMsg)).Else(validation.Required.Error(eitherMsg))),
		validation.Field(&f.Name, validation.When(f.ID != 0, validation.Empty.Error(bothMsg)).Else(validation.Required.Error(eitherMsg))),
		validation.Field(&f.Value, validation.Required),
	)
}

// Attachment allows to attach attachments to tasks and comments.
type Attachment struct {
	GUID         string `json:"guid,omitempty"`
	RootID       int    `json:"root_id,omitempty"`
	AttachmentID int    `json:"attachment_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
}

// Validate allows to validate request before sending.
func (a Attachment) Validate() error {
	const onlyOneMsg = "use guid, attachment_id or url, but not simultaneously"
	const eitherMsg = "use either guid, attachment or url"

	return validation.ValidateStruct(
		&a,
		validation.Field(&a.GUID,
			validation.When(
				a.AttachmentID != 0 || a.URL != "",
				validation.Empty.Error(onlyOneMsg),
			).Else(validation.Required.Error(eitherMsg)),
			validation.When(a.RootID != 0, validation.Required.Error("root_id requires guid")),
			is.UUID,
		),
		validation.Field(&a.AttachmentID, validation.When(
			a.GUID != "" || a.URL != "",
			validation.Empty.Error(onlyOneMsg),
		).Else(validation.Required.Error(eitherMsg))),
		validation.Field(&a.URL,
			validation.When(
				a.GUID != "" || a.AttachmentID != 0,
				validation.Empty.Error(onlyOneMsg),
			).Else(validation.Required.Error(eitherMsg)),
			validation.When(a.Name != "", validation.Required.Error("name requires url")),
			is.URL,
		),
	)
}

// TaskCommentRequest is necessary to create a comment in the task.
type TaskCommentRequest struct {
	Text                   string        `json:"text,omitempty"`
	Subject                string        `json:"subject,omitempty"`
	DueDate                string        `json:"due_date,omitempty"`
	Due                    *time.Time    `json:"due,omitempty"`
	Duration               int           `json:"duration,omitempty"`
	Action                 ActionType    `json:"action,omitempty"`
	ApprovalChoice         ChoiceType    `json:"approval_choice,omitempty"`
	ReassignTo             *Person       `json:"reassign_to,omitempty"`
	ApprovalsAdded         [][]*Person   `json:"approvals_added,omitempty"`
	ApprovalsRemoved       [][]*Person   `json:"approvals_removed,omitempty"`
	ApprovalsRerequested   [][]*Person   `json:"approvals_rerequested,omitempty"`
	SubscribersAdded       []*Person     `json:"subscribers_added,omitempty"`
	SubscribersRemoved     []*Person     `json:"subscribers_removed,omitempty"`
	SubscribersRerequested []*Person     `json:"subscribers_rerequested,omitempty"`
	ParticipantsAdded      []*Person     `json:"participants_added,omitempty"`
	ParticipantsRemoved    []*Person     `json:"participants_removed,omitempty"`
	FieldUpdates           []*FormField  `json:"field_updates,omitempty"`
	Attachments            []*Attachment `json:"attachments,omitempty"`
	AddedListIDs           []int         `json:"added_list_ids,omitempty"`
	RemovedListIDs         []int         `json:"removed_list_ids,omitempty"`
	ScheduledDate          string        `json:"scheduled_date,omitempty"`
	ScheduledDatetimeUTC   *time.Time    `json:"scheduled_datetime_utc,omitempty"`
	CancelSchedule         bool          `json:"cancel_schedule,omitempty"`
	Channel                *Channel      `json:"channel,omitempty"`
	SpentMinutes           int           `json:"spent_minutes,omitempty"`
}

// Validate allows to validate request before sending.
func (r TaskCommentRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Due,
			validation.When(r.DueDate != "", validation.Nil.Error("use due or due_date, not both"), validation.Date("2006-01-02")),
			validation.When(r.Duration != 0, validation.Required.Error("duration requires due")),
		),
		validation.Field(&r.DueDate, validation.When(r.Due != nil, validation.Empty.Error("use due or due_date, not both"))),
		validation.Field(&r.ReassignTo),
		validation.Field(&r.ApprovalsAdded, validation.Each()),
		validation.Field(&r.ApprovalsRemoved, validation.Each()),
		validation.Field(&r.ApprovalsRerequested, validation.Each()),
		validation.Field(&r.SubscribersAdded, validation.Each()),
		validation.Field(&r.SubscribersRemoved, validation.Each()),
		validation.Field(&r.SubscribersRerequested, validation.Each()),
		validation.Field(&r.ParticipantsAdded, validation.Each()),
		validation.Field(&r.ParticipantsRemoved, validation.Each()),
		validation.Field(&r.Attachments, validation.Each()),
		validation.Field(&r.ScheduledDate, validation.Date("2006-01-02")),
	)
}

// RegistryRequest is helpful to get a registry of tasks.
type RegistryRequest struct {
	FieldFilters map[int]string `json:"-"`

	Steps           int        `json:"steps,omitempty"`
	IncludeArchived bool       `json:"include_archived,omitempty"`
	FieldIDs        []int      `json:"field_ids,omitempty"`
	Format          string     `json:"format,omitempty"`
	Delimiter       string     `json:"delimiter,omitempty"`
	Encoding        string     `json:"encoding,omitempty"`
	SimpleFormat    bool       `json:"simple_format,omitempty"`
	ModifiedBefore  *time.Time `json:"modified_before,omitempty"`
	ModifiedAfter   *time.Time `json:"modified_after,omitempty"`
	CreatedBefore   *time.Time `json:"created_before,omitempty"`
	CreatedAfter    *time.Time `json:"created_after,omitempty"`
	ClosedBefore    *time.Time `json:"closed_before,omitempty"`
	ClosedAfter     *time.Time `json:"closed_after,omitempty"`
}

// MarshalJSON is a custom RegistryRequest marshaller that allows to merge the main struct and a map of field filters.
func (r *RegistryRequest) MarshalJSON() ([]byte, error) {
	if r.FieldFilters == nil {
		type Alias RegistryRequest
		aux := &struct {
			IncludeArchived string `json:"include_archived,omitempty"`
			SimpleFormat    string `json:"simple_format,omitempty"`
			*Alias
		}{
			Alias: (*Alias)(r),
		}
		if r.IncludeArchived {
			aux.IncludeArchived = "y"
		}
		if r.SimpleFormat {
			aux.SimpleFormat = "y"
		}
		return json.Marshal(aux)
	}

	m := make(map[string]interface{})
	for k, v := range r.FieldFilters {
		m["fld"+strconv.Itoa(k)] = v
	}

	t := reflect.TypeOf(r).Elem()
	v := reflect.ValueOf(r).Elem()
	for i := 0; i < t.NumField(); i++ {
		k := strings.ReplaceAll(t.Field(i).Tag.Get("json"), ",omitempty", "")
		switch k {
		case "-":
			continue
		case "include_archived":
			if v.Field(i).Interface().(bool) {
				m[k] = "y"
				continue
			}
		case "simple_format":
			if v.Field(i).Interface().(bool) {
				m[k] = "y"
				continue
			}
		}

		if !v.Field(i).IsZero() {
			m[k] = v.Field(i).Interface()
		}
	}

	return json.Marshal(m)
}

type fileRequest struct {
	Filename string
	io.Reader
}

type catalogRequest struct {
	Name           string         `json:"name"`
	CatalogHeaders []string       `json:"catalog_headers"`
	Items          []*CatalogItem `json:"items"`
}

type syncCatalogRequest struct {
	Apply          bool           `json:"apply"`
	CatalogHeaders []string       `json:"catalog_headers"`
	Items          []*CatalogItem `json:"items"`
}

// MemberRequest is necessary to create and update Member.
type MemberRequest struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Email        string `json:"email,omitempty"`
	Position     string `json:"position,omitempty"`
	DepartmentID int    `json:"department_id,omitempty"`
	Skype        string `json:"skype,omitempty"`
	Phone        string `json:"phone,omitempty"`
}

type roleRequest struct {
	Name      string `json:"name"`
	MemberAdd []int  `json:"member_add"`
}

type roleUpdateRequest struct {
	Name         string `json:"name,omitempty"`
	MemberAdd    []int  `json:"member_add,omitempty"`
	MemberRemove []int  `json:"member_remove,omitempty"`
	Banned       bool   `json:"banned"`
}

// RegisterCallRequest is necessary to register a call.
type RegisterCallRequest struct {
	To              string `json:"to,omitempty"`
	From            string `json:"from"`
	Extension       string `json:"extension,omitempty"`
	IntegrationGUID string `json:"integration_guid"`
	CallGUID        string `json:"call_guid,omitempty"`
	TaskID          int    `json:"task_id,omitempty"`
}

// Validate allows to validate request before sending.
func (r RegisterCallRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.From, validation.Required),
		validation.Field(&r.IntegrationGUID, validation.Required),
	)
}

// AddCallDetailsRequest is necessary to add call details.
type AddCallDetailsRequest struct {
	StartTime       *time.Time          `json:"start_time,omitempty"`
	EndTime         *time.Time          `json:"end_time,omitempty"`
	Rating          int                 `json:"rating,omitempty"`
	DisconnectParty DisconnectPartyType `json:"disconnect_party,omitempty"`
	CallStatus      CallStatusType      `json:"call_status,omitempty"`
	FileGUID        string              `json:"file_guid"`
}

type registerCallEventRequest struct {
	EventType CallEventType `json:"event_type"`
	Extension string        `json:"extension,omitempty"`
}
