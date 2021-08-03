package pyrus

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	baseURL   = "https://api.pyrus.com/v4"
	userAgent = "Pyrus API golang client v0.0.1"
)

type Client struct {
	baseURL string

	login       string
	securityKey string

	accessToken string
	mu          sync.RWMutex

	logger          Logger
	httpClient      *http.Client
	eventBufferSize int
}

// IClient is the main interface. Provided to implement dummy implementations useful for testing.
type IClient interface {
	Auth(login, securityKey string) (string, error)
	Forms() (*FormsResponse, error)
	Form(formID int) (*FormResponse, error)
	Registry(formID int, req *RegistryRequest) (*FormRegisterResponse, error)
	Task(taskID int) (*TaskResponse, error)
	CreateTask(req *TaskRequest) (*TaskResponse, error)
	CommentTask(taskID int, req *TaskCommentRequest) (*TaskResponse, error)
	UploadFile(name string, file io.Reader) (*UploadResponse, error)
	DownloadFile(fileID int) (*DownloadResponse, error)
	Catalogs() (*CatalogsResponse, error)
	Catalog(catalogID int) (*CatalogResponse, error)
	CreateCatalog(name string, headers []string, items []*CatalogItem) (*CatalogResponse, error)
	SyncCatalog(catalogID int, apply bool, headers []string, items []*CatalogItem) (*SyncCatalogResponse, error)
	Contacts() (*ContactsResponse, error)
	Members() (*MembersResponse, error)
	CreateMember(req *MemberRequest) (*Member, error)
	UpdateMember(memberID int, req *MemberRequest) (*Member, error)
	BlockMember(memberID int) (*Member, error)
	Roles() (*RolesResponse, error)
	CreateRole(name string, members []int) (*Role, error)
	UpdateRole(roleID int, name string, add, remove []int, banned bool) (*Role, error)
	Profile() (*ProfileResponse, error)
	Lists() (*ListsResponses, error)
	TaskList(listID, itemCount int, includeArchived bool) (*TaskListResponse, error)
	Inbox(itemCount int) (*TaskListResponse, error)
	RegisterCall(req *RegisterCallRequest) (*RegisterCallResponse, error)
	AddCallDetails(callGUID string, req *AddCallDetailsRequest) error
	RegisterCallEvent(callGUID string, eventType CallEventType, extension string) error
	WebhookHandler() (http.HandlerFunc, <-chan Event)
}

// Option helps to create an option for Client.
type Option func(*Client)

// WithLogger allows to log errors with own logger.
func WithLogger(l Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

// WithZapLogger allows to pass ready *zap.Logger instance for error logging.
func WithZapLogger(l *zap.Logger) Option {
	return func(c *Client) {
		c.logger = &zapLogger{logger: l}
	}
}

// WithHTTPClient allows to override http.DefaultClient and use your own.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithEventBufferSize allows to override default buffer size of Event chan used by Webhook engine.
func WithEventBufferSize(size int) Option {
	return func(c *Client) {
		c.eventBufferSize = size
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient returns an instance of Client.
func NewClient(login, securityKey string, opts ...Option) (*Client, error) {
	c := &Client{
		baseURL: baseURL,

		login:       login,
		securityKey: securityKey,

		logger:          &noopLogger{},
		httpClient:      http.DefaultClient,
		eventBufferSize: 100,
	}

	// Apply optional opts
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) getAndSetAccessToken() error {
	accessToken, err := c.Auth(c.login, c.securityKey)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.accessToken = accessToken
	c.mu.Unlock()

	return nil
}

func (c *Client) performRequest(method, path string, q *url.Values, reqBody, respBody interface{}) error {
	auth := false
	if path == "/auth" {
		auth = true
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		c.logger.Error("Error while parsing a URL!", err)
		return err
	}
	if q != nil {
		u.RawQuery = q.Encode()
	}

	multipartRequest := false
	if _, ok := reqBody.(*fileRequest); ok {
		multipartRequest = true
	}

	var (
		req    *http.Request
		reqErr error
	)
	contentTypeHeader := "application/json"
	if multipartRequest {
		buf := bytes.NewBuffer(nil)

		w := multipart.NewWriter(buf)
		fw, err := w.CreateFormFile("file", reqBody.(*fileRequest).Filename)
		if err != nil {
			c.logger.Error("Error while creating a new form file!", err)
			return err
		}
		if _, err := io.Copy(fw, reqBody.(*fileRequest).Reader); err != nil {
			c.logger.Error("Error while writing a file!", err)
			return err
		}
		if err := w.Close(); err != nil {
			c.logger.Error("Error while trying to close multipart writer!", err)
			return err
		}

		req, reqErr = http.NewRequest(method, u.String(), buf)
		contentTypeHeader = w.FormDataContentType()
	} else if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
			c.logger.Error("Error while encoding JSON!", err)
			return err
		}

		req, reqErr = http.NewRequest(method, u.String(), buf)
	} else {
		req, reqErr = http.NewRequest(method, u.String(), nil)
	}
	if reqErr != nil {
		c.logger.Error("Error while creating a request!", err)
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentTypeHeader)

	// It's wise to get first token without unnecessary request
	c.mu.RLock()
	ok := c.accessToken != "" || auth
	c.mu.RUnlock()
	if !ok {
		if err := c.getAndSetAccessToken(); err != nil {
			return err
		}
	}

	c.mu.RLock()
	if c.accessToken != "" && !auth {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error while doing a request!", err)
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	// Get new access_token in case of old session
	if resp.StatusCode == 401 && !auth {
		if err := c.getAndSetAccessToken(); err != nil {
			return err
		}

		return c.performRequest(method, path, q, reqBody, respBody)
	}

	// Don't read if there is no need in response body at all
	if respBody == nil && !auth {
		return nil
	}

	// File downloading
	if mt, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type")); err == nil && mt != "application/json" {
		mt, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err != nil {
			c.logger.Error("Error while parsing media type!", err)
			return err
		}
		if mt != "attachment" {
			return errors.New("attachment was expected")
		}

		filename, ok := params["filename"]
		if !ok {
			return errors.New("file doesn't have a name")
		}
		if _, ok := respBody.(*string); ok {
			*respBody.(*string) = filename
		}

		w, ok := reqBody.(io.Writer)
		if !ok {
			return errors.New("writer was expected")
		}

		if _, err := io.Copy(w, resp.Body); err != nil {
			c.logger.Error("Error while trying to download file!", err)
			return err
		}

		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var pe Error
		if err := decoder.Decode(&pe); err != nil {
			c.logger.Error("Error while decoding a response body!", err)
			return err
		}

		return pe
	}

	if err := decoder.Decode(&respBody); err != nil {
		c.logger.Error("Error while decoding a response body!", err)
		return err
	}

	return nil
}

// Auth performs authorization and returns access_token.
func (c *Client) Auth(login, securityKey string) (string, error) {
	var respBody AuthResponse
	if err := c.performRequest(http.MethodPost, "/auth", nil, &authRequest{
		Login:       login,
		SecurityKey: securityKey,
	}, &respBody); err != nil {
		return "", err
	}

	return respBody.AccessToken, nil
}

// Forms returns a description of all the forms in which the current user is a manager or a member.
func (c *Client) Forms() (*FormsResponse, error) {
	var forms FormsResponse
	if err := c.performRequest(http.MethodGet, "/forms", nil, nil, &forms); err != nil {
		return nil, err
	}

	return &forms, nil
}

// Form returns a description of form with inputted id.
func (c *Client) Form(formID int) (*FormResponse, error) {
	var form FormResponse
	if err := c.performRequest(http.MethodGet, "/forms/"+strconv.Itoa(formID), nil, nil, &form); err != nil {
		return nil, err
	}

	return &form, nil
}

// Registry returns the list of tasks that were created based on the specified form.
// The response only contains general information about the task, like the list of filled form fields and its workflow.
// You can use Task method to get all task comments.
func (c *Client) Registry(formID int, req *RegistryRequest) (*FormRegisterResponse, error) {
	var tasks FormRegisterResponse
	if err := c.performRequest(http.MethodPost, "/forms/"+strconv.Itoa(formID)+"/register", nil, req, &tasks); err != nil {
		return nil, err
	}

	return &tasks, nil
}

// Task returns a task with all comments.
func (c *Client) Task(taskID int) (*TaskResponse, error) {
	var task TaskResponse
	if err := c.performRequest(http.MethodGet, "/tasks/"+strconv.Itoa(taskID), nil, nil, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CreateTask creates a task and returns it with a comment.
func (c *Client) CreateTask(req *TaskRequest) (*TaskResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var task TaskResponse
	if err := c.performRequest(http.MethodPost, "/tasks", nil, req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CommentTask comments a task and returns it with all comments, including the added one.
func (c *Client) CommentTask(taskID int, req *TaskCommentRequest) (*TaskResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var task TaskResponse
	if err := c.performRequest(http.MethodPost, "/tasks/"+strconv.Itoa(taskID)+"/comments", nil, req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UploadFile uploads files for subsequent attachment to tasks.
// Files that are not referenced by any task are removed after a while.
func (c *Client) UploadFile(name string, file io.Reader) (*UploadResponse, error) {
	var upload UploadResponse
	if err := c.performRequest(http.MethodPost, "/files/upload", nil, &fileRequest{
		Filename: name,
		Reader:   file,
	}, &upload); err != nil {
		return nil, err
	}

	return &upload, nil
}

// DownloadFile downloads file from Pyrus.
func (c *Client) DownloadFile(fileID int) (*DownloadResponse, error) {
	buf := bytes.NewBuffer(nil)

	var filename string
	if err := c.performRequest(http.MethodGet, "/files/download/"+strconv.Itoa(fileID), nil, buf, &filename); err != nil {
		return nil, err
	}

	return &DownloadResponse{
		Filename: filename,
		RawFile:  buf.Bytes(),
	}, nil
}

// Catalogs returns a list of available catalogs.
func (c *Client) Catalogs() (*CatalogsResponse, error) {
	var catalogs CatalogsResponse
	if err := c.performRequest(http.MethodGet, "/catalogs", nil, nil, &catalogs); err != nil {
		return nil, err
	}

	return &catalogs, nil
}

// Catalog returns a catalog with all its elements.
func (c *Client) Catalog(catalogID int) (*CatalogResponse, error) {
	var catalog CatalogResponse
	if err := c.performRequest(http.MethodGet, "/catalogs/"+strconv.Itoa(catalogID), nil, nil, &catalog); err != nil {
		return nil, err
	}

	return &catalog, nil
}

// CreateCatalog creates a catalog and returns it with all its elements.
func (c *Client) CreateCatalog(name string, headers []string, items []*CatalogItem) (*CatalogResponse, error) {
	var catalog CatalogResponse
	if err := c.performRequest(http.MethodPut, "/catalogs", nil, &catalogRequest{
		Name:           name,
		CatalogHeaders: headers,
		Items:          items,
	}, &catalog); err != nil {
		return nil, err
	}

	return &catalog, nil
}

// SyncCatalog updates catalog header and items and returns a list of items that have been added, modified, or deleted.
func (c *Client) SyncCatalog(catalogID int, apply bool, headers []string, items []*CatalogItem) (*SyncCatalogResponse, error) {
	var syncCatalog SyncCatalogResponse
	if err := c.performRequest(http.MethodPost, "/catalogs/"+strconv.Itoa(catalogID), nil, &syncCatalogRequest{
		Apply:          apply,
		CatalogHeaders: headers,
		Items:          items,
	}, &syncCatalog); err != nil {
		return nil, err
	}

	return &syncCatalog, nil
}

// Contacts returns a list of contacts available to the current user and grouped by organization.
func (c *Client) Contacts() (*ContactsResponse, error) {
	var contacts ContactsResponse
	if err := c.performRequest(http.MethodGet, "/contacts", nil, nil, &contacts); err != nil {
		return nil, err
	}

	return &contacts, nil
}

// Members returns a list of all organization participants.
func (c *Client) Members() (*MembersResponse, error) {
	var members MembersResponse
	if err := c.performRequest(http.MethodGet, "/members", nil, nil, &members); err != nil {
		return nil, err
	}

	return &members, nil
}

// CreateMember creates a user and returns it.
func (c *Client) CreateMember(req *MemberRequest) (*Member, error) {
	var member Member
	if err := c.performRequest(http.MethodPost, "/members", nil, req, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// UpdateMember updates a user and returns it.
func (c *Client) UpdateMember(memberID int, req *MemberRequest) (*Member, error) {
	var member Member
	if err := c.performRequest(http.MethodPut, "/members/"+strconv.Itoa(memberID), nil, req, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// BlockMember blocks a user and returns it.
func (c *Client) BlockMember(memberID int) (*Member, error) {
	var member Member
	if err := c.performRequest(http.MethodDelete, "/members/"+strconv.Itoa(memberID), nil, nil, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// Roles returns a list of roles.
func (c *Client) Roles() (*RolesResponse, error) {
	var roles RolesResponse
	if err := c.performRequest(http.MethodGet, "/roles", nil, nil, &roles); err != nil {
		return nil, err
	}

	return &roles, nil
}

// CreateRole creates a role and returns it.
func (c *Client) CreateRole(name string, members []int) (*Role, error) {
	var role Role
	if err := c.performRequest(http.MethodPost, "/roles", nil, &roleRequest{
		Name:      name,
		MemberAdd: members,
	}, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRole updates a role and returns it.
func (c *Client) UpdateRole(roleID int, name string, add, remove []int, banned bool) (*Role, error) {
	var role Role
	if err := c.performRequest(http.MethodPut, "/roles/"+strconv.Itoa(roleID), nil, &roleUpdateRequest{
		Name:         name,
		MemberAdd:    add,
		MemberRemove: remove,
		Banned:       banned,
	}, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// Profile returns a profile of the calling user.
func (c *Client) Profile() (*ProfileResponse, error) {
	var profile ProfileResponse
	if err := c.performRequest(http.MethodGet, "/profile", nil, nil, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Lists returns all the lists that are available to the user.
func (c *Client) Lists() (*ListsResponses, error) {
	var lists ListsResponses
	if err := c.performRequest(http.MethodGet, "/lists", nil, nil, &lists); err != nil {
		return nil, err
	}

	return &lists, nil
}

// TaskList returns all the tasks in the specified list.
func (c *Client) TaskList(listID, itemCount int, includeArchived bool) (*TaskListResponse, error) {
	q := &url.Values{}
	if itemCount != 0 {
		q.Set("item_count", strconv.Itoa(itemCount))
	}
	if includeArchived {
		q.Set("include_archived", "y")
	}

	var taskList TaskListResponse
	if err := c.performRequest(http.MethodGet, "/lists/"+strconv.Itoa(listID)+"/tasks", q, nil, &taskList); err != nil {
		return nil, err
	}

	return &taskList, nil
}

// Inbox returns all inbox tasks.
func (c *Client) Inbox(itemCount int) (*TaskListResponse, error) {
	q := &url.Values{}
	if itemCount != 0 {
		q.Set("item_count", strconv.Itoa(itemCount))
	}

	var taskList TaskListResponse
	if err := c.performRequest(http.MethodGet, "/inbox", q, nil, &taskList); err != nil {
		return nil, err
	}

	return &taskList, nil
}

// RegisterCall returns the GUID of the incoming call, and the id of the generated request.
func (c *Client) RegisterCall(req *RegisterCallRequest) (*RegisterCallResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var call RegisterCallResponse
	if err := c.performRequest(http.MethodPost, "/calls", nil, req, &call); err != nil {
		return nil, err
	}

	return &call, nil
}

// AddCallDetails adds call details by call_guid.
func (c *Client) AddCallDetails(callGUID string, req *AddCallDetailsRequest) error {
	if err := c.performRequest(http.MethodPut, "/calls/"+callGUID, nil, req, nil); err != nil {
		return err
	}

	return nil
}

// RegisterCallEvent registers call event by call_guid.
func (c *Client) RegisterCallEvent(callGUID string, eventType CallEventType, extension string) error {
	if err := c.performRequest(http.MethodPost, "/calls/"+callGUID+"/event", nil, &registerCallEventRequest{
		EventType: eventType,
		Extension: extension,
	}, nil); err != nil {
		return err
	}

	return nil
}

// WebhookHandler returns HTTP handler and channel with Event's.
// Handler automatically checks X-Pyrus-Sig, parses Event and sends it over channel..
func (c *Client) WebhookHandler() (http.HandlerFunc, <-chan Event) {
	eventChan := make(chan Event, c.eventBufferSize)

	writeError := func(w http.ResponseWriter, code int, err error) {
		respBody, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(respBody); err != nil {
			c.logger.Error("Error while writing a response!", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			c.logger.Error("Error while reading a request body!", err)
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		hasher := hmac.New(sha1.New, []byte(c.securityKey))
		hasher.Write(b)
		hash := hex.EncodeToString(hasher.Sum(nil))
		if subtle.ConstantTimeCompare([]byte(hash), []byte(strings.ToLower(r.Header.Get("X-Pyrus-Sig")))) != 1 {
			err := errors.New("invalid signature")
			c.logger.Error("Invalid signature!", err)
			writeError(w, http.StatusUnauthorized, err)
			return
		}

		var event Event
		if err := json.Unmarshal(b, &event); err != nil {
			c.logger.Error("Error while decoding a request body!", err)
			writeError(w, http.StatusBadRequest, err)
			return
		}

		eventChan <- event
		w.WriteHeader(http.StatusOK)
	}, eventChan
}
