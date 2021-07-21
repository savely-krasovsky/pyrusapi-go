package pyrus

import (
	"encoding/json"
	"time"
)

// FormField is a Form field. Forms consist of fields.
// They could usually have tree structure, so often you will have to use type assertion.
type FormField struct {
	ID    int            `json:"id"`
	Type  FieldType      `json:"type"`
	Name  string         `json:"name"`
	Info  *FormFieldInfo `json:"info,omitempty"`
	Value interface{}    `json:"value,omitempty"`
	// ParentID returns if field has parent
	ParentID int `json:"parent_id,omitempty"`
	// RowID returns if field is in table
	RowID int `json:"row_id,omitempty"`
}

// FormFieldInfo could contain additional field information
type FormFieldInfo struct {
	// RequiredStep indicates a step number where a field becomes required for filling
	RequiredStep int `json:"required_step"`
	// ImmutableStep indicates a step number from which the user can't change a field value
	ImmutableStep int `json:"immutable_step"`
	// Options return for a multiple_choice field
	Options []*ChoiceOption `json:"options,omitempty"`
	// CatalogID returns for a catalog field
	CatalogID int `json:"catalog_id,omitempty"`
	// Columns return for a table field
	Columns []*FormField `json:"columns,omitempty"`
	// Fields return for a title field
	Fields []*FormField `json:"fields,omitempty"`
	// DecimalPlaces return for a number field
	DecimalPlaces int `json:"decimal_places,omitempty"`
}

// ChoiceOption represents a choice option of multiple_choice field type.
type ChoiceOption struct {
	ChoiceID    int          `json:"choice_id"`
	ChoiceValue string       `json:"choice_value"`
	Fields      []*FormField `json:"fields"`
	Deleted     bool         `json:"deleted"`
}

// UnmarshalJSON is a custom unmarshaler to create a tree of form fields.
func (f *FormField) UnmarshalJSON(b []byte) error {
	type RawFormField FormField
	raw := &struct {
		Value json.RawMessage `json:"value"`
		*RawFormField
	}{
		RawFormField: (*RawFormField)(f),
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw.Value == nil {
		return nil
	}

	switch raw.Type {
	case FieldTypeText:
		var text string
		if err := json.Unmarshal(raw.Value, &text); err != nil {
			return err
		}

		f.Value = text
	case FieldTypeMoney:
		var money float64
		if err := json.Unmarshal(raw.Value, &money); err != nil {
			return err
		}

		f.Value = money
	case FieldTypeNumber:
		var number float64
		if err := json.Unmarshal(raw.Value, &number); err != nil {
			return err
		}

		f.Value = number
	case FieldTypeDate:
		var dateStr string
		if err := json.Unmarshal(raw.Value, &dateStr); err != nil {
			return err
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		f.Value = date
	case FieldTypeTime:
		var timeStr string
		if err := json.Unmarshal(raw.Value, &timeStr); err != nil {
			return err
		}

		t, err := time.Parse("15:04", timeStr)
		if err != nil {
			return err
		}

		f.Value = t
	case FieldTypeCheckmark:
		var checkmark CheckmarkType
		if err := json.Unmarshal(raw.Value, &checkmark); err != nil {
			return err
		}

		f.Value = checkmark
	case FieldTypeDueDate:
		var dateStr string
		if err := json.Unmarshal(raw.Value, &dateStr); err != nil {
			return err
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		f.Value = date
	case FieldTypeDueDateTime:
		var dateStr string
		if err := json.Unmarshal(raw.Value, &dateStr); err != nil {
			return err
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return err
		}

		f.Value = date
	case FieldTypeEmail:
		var email string
		if err := json.Unmarshal(raw.Value, &email); err != nil {
			return err
		}

		f.Value = email
	case FieldTypePhone:
		var phone string
		if err := json.Unmarshal(raw.Value, &phone); err != nil {
			return err
		}

		f.Value = phone
	case FieldTypeFlag:
		var flg FlagType
		if err := json.Unmarshal(raw.Value, &flg); err != nil {
			return err
		}

		f.Value = flg
	case FieldTypeStep:
		var step int
		if err := json.Unmarshal(raw.Value, &step); err != nil {
			return err
		}

		f.Value = step
	case FieldTypeStatus:
		var status StatusType
		if err := json.Unmarshal(raw.Value, &status); err != nil {
			return err
		}

		f.Value = status
	case FieldTypeCreationDate:
		var dateStr string
		if err := json.Unmarshal(raw.Value, &dateStr); err != nil {
			return err
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		f.Value = date
	case FieldTypeNote:
		var note string
		if err := json.Unmarshal(raw.Value, &note); err != nil {
			return err
		}

		f.Value = note
	case FieldTypeCatalog:
		var catalogItem CatalogItem
		if err := json.Unmarshal(raw.Value, &catalogItem); err != nil {
			return err
		}

		f.Value = &catalogItem
	case FieldTypeFile:
		var files []*File
		if err := json.Unmarshal(raw.Value, &files); err != nil {
			return err
		}

		f.Value = files
	case FieldTypePerson:
		var person Person
		if err := json.Unmarshal(raw.Value, &person); err != nil {
			return err
		}

		f.Value = &person
	case FieldTypeAuthor:
		var author Person
		if err := json.Unmarshal(raw.Value, &author); err != nil {
			return err
		}

		f.Value = &author
	case FieldTypeTable:
		var table Table
		if err := json.Unmarshal(raw.Value, &table); err != nil {
			return err
		}

		f.Value = table
	case FieldTypeMultipleChoice:
		var mc MultipleChoice
		if err := json.Unmarshal(raw.Value, &mc); err != nil {
			return err
		}

		f.Value = &mc
	case FieldTypeTitle:
		var title Title
		if err := json.Unmarshal(raw.Value, &title); err != nil {
			return err
		}

		f.Value = &title
	case FieldTypeFormLink:
		var formLink FormLink
		if err := json.Unmarshal(raw.Value, &formLink); err != nil {
			return err
		}

		f.Value = &formLink
	default:
		var i interface{}
		if err := json.Unmarshal(raw.Value, &i); err != nil {
			return err
		}

		f.Value = i
	}

	return nil
}

// TaskHeader represents only basic information about a task.
type TaskHeader struct {
	ID               int        `json:"id"`
	CreateDate       time.Time  `json:"create_date"`
	LastModifiedDate *time.Time `json:"last_modified_date"`
	CloseDate        *time.Time `json:"close_date"`
	Author           *Person    `json:"author"`

	Text        string  `json:"text"`
	Responsible *Person `json:"responsible"`
}

// Task represents a task without comments.
type Task struct {
	*TaskHeader

	Attachments          []*File       `json:"attachments"`
	ListIDs              []int         `json:"list_ids"`
	ParentTaskID         int           `json:"parent_task_id"`
	LinkedTaskIDs        []int         `json:"linked_task_ids"`
	LastNoteID           int           `json:"last_note_id"`
	Subject              string        `json:"subject"`
	ScheduledDate        string        `json:"scheduled_date"`
	ScheduledDatetimeUTC *time.Time    `json:"scheduled_datetime_utc"`
	Subscribers          []*Subscriber `json:"subscribers"`

	DueDate      string     `json:"due_date"`
	Due          *time.Time `json:"due"`
	Duration     int        `json:"duration"`
	Participants []*Person  `json:"participants"`

	FormID      int           `json:"form_id"`
	Fields      []*FormField  `json:"fields,omitempty"`
	Approvals   [][]*Approval `json:"approvals"`
	CurrentStep int           `json:"current_step"`
}

// TaskWithComments represents a task with all of its comments.
type TaskWithComments struct {
	*Task

	Comments []*TaskComment `json:"comments,omitempty"`
}

// Person represents a user of Pyrus.
type Person struct {
	ID             int        `json:"id,omitempty"`
	FirstName      string     `json:"first_name,omitempty"`
	LastName       string     `json:"last_name,omitempty"`
	Email          string     `json:"email,omitempty"`
	Type           PersonType `json:"type,omitempty"`
	DepartmentID   int        `json:"department_id,omitempty"`
	DepartmentName string     `json:"department_name,omitempty"`
}

// File represents an attachment to the task. It could be a part of filled form or comment.
type File struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Size    int    `json:"size"`
	MD5     string `json:"md5"`
	URL     string `json:"url"`
	Version int    `json:"version"`
	RootID  int    `json:"root_id"`
}

// Approval represents an approval by person. It contains person, step number and choice itself.
type Approval struct {
	Person         *Person    `json:"person"`
	Step           int        `json:"step"`
	ApprovalChoice ChoiceType `json:"approval_choice"`
}

// Subscriber represents a person who can watch for task updates, but doesn't participate in the process of approval.
type Subscriber struct {
	Person         *Person    `json:"person"`
	ApprovalChoice ChoiceType `json:"approval_choice"`
}

// TaskComment represents a comment from task. Comment is not only the text,
// it contains all the updates of tasks: field updates, approvals, reassignments, etc.
type TaskComment struct {
	ID                     int        `json:"id"`
	Text                   string     `json:"text"`
	CreateDate             time.Time  `json:"create_date"`
	Author                 *Person    `json:"author"`
	Attachments            []*File    `json:"attachments"`
	Action                 ActionType `json:"action"`
	AddedListIDs           []int      `json:"added_list_ids"`
	RemovedListIDs         []int      `json:"removed_list_ids"`
	CommentAsRoles         []*Role    `json:"comment_as_roles"`
	Subject                string     `json:"subject"`
	ScheduledDate          string     `json:"scheduled_date"`
	ScheduledDatetimeUTC   *time.Time `json:"scheduled_datetime_utc"`
	CancelSchedule         bool       `json:"cancel_schedule"`
	SpentMinutes           int        `json:"spent_minutes"`
	SubscribersAdded       []*Person  `json:"subscribers_added"`
	SubscribersRemoved     []*Person  `json:"subscribers_removed"`
	SubscribersRerequested []*Person  `json:"subscribers_rerequested"`

	ReassignedTo        *Person    `json:"reassigned_to"`
	ParticipantsAdded   []*Person  `json:"participants_added"`
	ParticipantsRemoved []*Person  `json:"participants_removed"`
	DueDate             string     `json:"due_date"`
	Due                 *time.Time `json:"due"`
	Duration            int        `json:"duration"`

	FieldUpdates         []*FormField  `json:"field_updates"`
	ApprovalChoice       ChoiceType    `json:"approval_choice"`
	ApprovalStep         int           `json:"approval_step"`
	ResetToStep          int           `json:"reset_to_step"`
	ChangedStep          int           `json:"changed_step"`
	ApprovalsAdded       [][]*Approval `json:"approvals_added"`
	ApprovalsRemoved     [][]*Approval `json:"approvals_removed"`
	ApprovalsRerequested [][]*Approval `json:"approvals_rerequested"`
	Channel              *Channel      `json:"channel"`
}

// Organization represents organization with persons and roles of it.
type Organization struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Persons []*Person `json:"persons"`
	Roles   []*Role   `json:"roles"`
}

// Role represents role and its members.
type Role struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MemberIDs  []int  `json:"member_ids"`
	ExternalID int    `json:"external_id"`
	Banned     bool   `json:"banned"`
}

// CatalogItem represents an item of Catalog. It contains headers of catalog and its value.
type CatalogItem struct {
	ItemID  int        `json:"item_id,omitempty"`
	ItemIDs []int      `json:"item_ids,omitempty"`
	Headers []string   `json:"headers,omitempty"`
	Values  []string   `json:"values,omitempty"`
	Rows    [][]string `json:"rows,omitempty"`
}

// CatalogHeader represents a header of Catalog. For example column "Name" or "Email".
type CatalogHeader struct {
	Name string            `json:"name"`
	Type CatalogHeaderType `json:"type"`
}

// Table represents a table. In our case it's just a slice of table rows.
type Table []*TableRow

// TableRow is an element of table.
type TableRow struct {
	RowID  int          `json:"row_id"`
	Cells  []*FormField `json:"cells"`
	Delete bool         `json:"delete"`
}

// Title represents a form field title (official docs doesn't explain what exactly it is).
type Title struct {
	Checkmark CheckmarkType `json:"checkmark"`
	Fields    []*FormField  `json:"fields"`
}

// MultipleChoice represents a form field with multiple choice dropdown menu.
type MultipleChoice struct {
	ChoiceIDs   []int        `json:"choice_ids"`
	ChoiceNames []string     `json:"choice_names"`
	Fields      []*FormField `json:"fields"`
	ChoiceID    int          `json:"choice_id"`
}

// FormLink represents a form field (official docs doesn't explain what exactly it is).
type FormLink struct {
	TaskIDs []int  `json:"task_ids"`
	Subject string `json:"subject"`
}

// Channel represents an external channel of comments. It allows to mark there to send or from there it was sent.
type Channel struct {
	Type ChannelType  `json:"type"`
	To   *ChannelUser `json:"to"`
	From *ChannelUser `json:"from"`
}

// ChannelUser represents a user from Channel. Email is used only for email and Name for everything else.
type ChannelUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TaskList represents a list of tasks.
type TaskList struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Children []*TaskList `json:"children"`
}

// Member represents a member of organization.
type Member struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Type      PersonType `json:"type"`
	Banned    bool       `json:"banned"`
	Position  string     `json:"position"`
	Skype     string     `json:"skype"`
	Phone     string     `json:"phone"`
}
