package pyrus

// AuthResponse represents a response from Auth method.
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// FormResponse represents a response from Form method.
type FormResponse struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Steps           map[int]string `json:"steps"`
	Fields          []*FormField   `json:"fields"`
	DeletedOrClosed bool           `json:"deleted_or_closed"`
	PrintForms      []PrintForm    `json:"print_forms"`
	Folder          []string       `json:"folder"`
}

type PrintForm struct {
	ID   int    `json:"print_form_id"`
	Name string `json:"print_form_name"`
}

// FormsResponse represents a response from Forms method.
type FormsResponse struct {
	Forms []*FormResponse `json:"forms"`
}

// FormRegisterResponse represents a response from Registry method.
type FormRegisterResponse struct {
	Tasks []*Task `json:"tasks"`
	CSV   string  `json:"csv"`
}

// TaskResponse represents a response from Task method.
type TaskResponse struct {
	Task *TaskWithComments `json:"task"`
}

// ContactsResponse represents a response from Contacts method.
type ContactsResponse struct {
	Organizations []*Organization `json:"organizations"`
}

// CatalogsResponse represents a list of available catalogs
type CatalogsResponse struct {
	Catalogs []*CatalogResponse `json:"catalogs"`
}

// CatalogResponse represents a response from Catalog method.
type CatalogResponse struct {
	CatalogID       int              `json:"catalog_id"`
	Name            string           `json:"name"`
	Version         int              `json:"version"`
	Supervisors     []int            `json:"supervisors"`
	Deleted         bool             `json:"deleted"`
	ExternalVersion int              `json:"external_version"`
	CatalogHeaders  []*CatalogHeader `json:"catalog_headers"`
	Items           []*CatalogItem   `json:"items"`
}

// UploadResponse represents a response from UploadFile method.
type UploadResponse struct {
	GUID    string `json:"guid"`
	MD5Hash string `json:"md5_hash"`
}

// DownloadResponse represents a response from DownloadFile method.
type DownloadResponse struct {
	Filename string
	RawFile  []byte `json:"raw_file"`
}

// ListsResponses represents a response from Lists method.
type ListsResponses struct {
	Lists []*TaskList `json:"lists"`
}

// TaskListResponse represents a response from TaskList method.
type TaskListResponse struct {
	Tasks   []*TaskHeader `json:"tasks"`
	HasMode bool          `json:"has_mode"`
}

// SyncCatalogResponse represents a response from SyncCatalog method.
type SyncCatalogResponse struct {
	Apply          bool             `json:"apply"`
	Added          []*CatalogItem   `json:"added"`
	Deleted        []*CatalogItem   `json:"deleted"`
	Updated        []*CatalogItem   `json:"updated"`
	CatalogHeaders []*CatalogHeader `json:"catalog_headers"`
}

// MembersResponse represents a response from Members method.
type MembersResponse struct {
	Members []*Member `json:"members"`
}

// RolesResponse represents a response from Roles method.
type RolesResponse struct {
	Roles []*Role `json:"roles"`
}

// ProfileResponse represents a response from Profile method.
type ProfileResponse struct {
	PersonID       int    `json:"person_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Locale         string `json:"locale"`
	OrganizationID int    `json:"organization_id"`
}

// RegisterCallResponse represents a response from RegisterCall method.
type RegisterCallResponse struct {
	CallGUID string `json:"call_guid"`
	TaskID   string `json:"task_id"`
}

// Event represents an event received from webhook.
type Event struct {
	Event       string            `json:"event"`
	AccessToken string            `json:"access_token"`
	TaskID      int               `json:"task_id"`
	UserID      int               `json:"user_id"`
	Task        *TaskWithComments `json:"task"`
}
