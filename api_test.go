package pyrus

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	pyrusLogin       string
	pyrusSecurityKey string
	formID           int
	taskID           int
	catalogID        int
	memberID         int
	roleID           int
	listID           int
	fileID           int
	callGUID         string

	logger, _ = zap.NewDevelopment()
	cl        Client
	ts        *httptest.Server
)

func TestMain(m *testing.M) {
	var (
		err    error
		exists bool
	)
	pyrusLogin, exists = os.LookupEnv("PYRUS_LOGIN")
	if !exists {
		log.Fatalln("Please pass PYRUS_LOGIN environment variable!")
	}
	pyrusSecurityKey, exists = os.LookupEnv("PYRUS_SECURITY_KEY")
	if !exists {
		log.Fatalln("Please pass PYRUS_SECURITY_KEY environment variable!")
	}
	formID, err = strconv.Atoi(os.Getenv("PYRUS_FORM_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_FORM_ID environment variable!")
	}
	taskID, err = strconv.Atoi(os.Getenv("PYRUS_TASK_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_TASK_ID environment variable!")
	}
	catalogID, err = strconv.Atoi(os.Getenv("PYRUS_CATALOG_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_CATALOG_ID environment variable!")
	}
	memberID, err = strconv.Atoi(os.Getenv("PYRUS_MEMBER_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_MEMBER_ID environment variable!")
	}
	roleID, err = strconv.Atoi(os.Getenv("PYRUS_ROLE_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_ROLE_ID environment variable!")
	}
	listID, err = strconv.Atoi(os.Getenv("PYRUS_LIST_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_LIST_ID environment variable!")
	}
	fileID, err = strconv.Atoi(os.Getenv("PYRUS_FILE_ID"))
	if err != nil {
		log.Fatalln("Please pass valid PYRUS_FILE_ID environment variable!")
	}
	callGUID, exists = os.LookupEnv("PYRUS_CALL_GUID")
	if !exists {
		log.Fatalln("Please pass PYRUS_CALL_GUID environment variable!")
	}

	var (
		requestAuth              = "POST:/auth"
		requestForms             = "GET:/forms"
		requestForm              = "GET:/forms/" + strconv.Itoa(formID)
		requestRegistry          = "POST:/forms/" + strconv.Itoa(formID) + "/register"
		requestTask              = "GET:/tasks/" + strconv.Itoa(taskID)
		requestCreateTask        = "POST:/tasks"
		requestCommentTask       = "POST:/tasks/" + strconv.Itoa(taskID) + "/comments"
		requestUploadFile        = "POST:/files/upload"
		requestDownloadFile      = "GET:/files/download/" + strconv.Itoa(fileID)
		requestCatalogs          = "GET:/catalogs"
		requestCatalog           = "GET:/catalogs/" + strconv.Itoa(catalogID)
		requestCreateCatalog     = "PUT:/catalogs"
		requestSyncCatalog       = "POST:/catalogs/" + strconv.Itoa(catalogID)
		requestContacts          = "GET:/contacts"
		requestMembers           = "GET:/members"
		requestCreateMember      = "POST:/members"
		requestUpdateMember      = "PUT:/members/" + strconv.Itoa(memberID)
		requestDeleteMember      = "DELETE:/members/" + strconv.Itoa(memberID)
		requestRoles             = "GET:/roles"
		requestCreateRole        = "POST:/roles"
		requestUpdateRole        = "PUT:/roles/" + strconv.Itoa(roleID)
		requestProfile           = "GET:/profile"
		requestLists             = "GET:/lists"
		requestListsTasks        = "GET:/lists/" + strconv.Itoa(listID) + "/tasks"
		requestInbox             = "GET:/inbox"
		requestRegisterCall      = "POST:/calls"
		requestAddCallDetails    = "PUT:/calls/" + callGUID
		requestRegisterCallEvent = "POST:/calls/" + callGUID + "/event"
	)

	requests := map[string]string{
		requestAuth:              "testdata/auth.json",
		requestForms:             "testdata/forms.json",
		requestForm:              "testdata/form.json",
		requestRegistry:          "testdata/registry.json",
		requestTask:              "testdata/task.json",
		requestCreateTask:        "testdata/task.json",
		requestCommentTask:       "testdata/task.json",
		requestUploadFile:        "testdata/uploaded_file.json",
		requestDownloadFile:      "testdata/downloaded_file.bin",
		requestCatalogs:          "testdata/catalogs.json",
		requestCatalog:           "testdata/catalog.json",
		requestCreateCatalog:     "testdata/catalog.json",
		requestSyncCatalog:       "testdata/sync_catalog.json",
		requestContacts:          "testdata/contacts.json",
		requestMembers:           "testdata/members.json",
		requestCreateMember:      "testdata/member.json",
		requestUpdateMember:      "testdata/member.json",
		requestDeleteMember:      "testdata/member.json",
		requestRoles:             "testdata/roles.json",
		requestCreateRole:        "testdata/role.json",
		requestUpdateRole:        "testdata/role.json",
		requestProfile:           "testdata/profile.json",
		requestLists:             "testdata/lists.json",
		requestListsTasks:        "testdata/lists_tasks.json",
		requestInbox:             "testdata/inbox.json",
		requestRegisterCall:      "testdata/call.json",
		requestAddCallDetails:    "",
		requestRegisterCallEvent: "",
	}

	// seed responses
	if pyrusLogin != "" && pyrusSecurityKey != "" {
		c, err := NewClient(pyrusLogin, pyrusSecurityKey)
		if err != nil {
			log.Fatalln(err)
		}

		accessToken, err := c.Auth(pyrusLogin, pyrusSecurityKey)
		if err != nil {
			log.Fatalln(err)
		}

		// forms
		for k, v := range requests {
			func() {
				mp := strings.Split(k, ":")
				method := mp[0]
				path := mp[1]

				if v == "" {
					return
				}
				if _, err := os.Stat(v); !errors.Is(err, fs.ErrNotExist) {
					return
				}

				// skip to receive data only from GETs
				switch {
				case method == "PUT" || method == "DELETE":
					return
				case method == "POST" && k != requestRegistry:
					return
				}

				f, err := os.Create(v)
				if err != nil {
					log.Println(err)
					return
				}

				defer f.Close() //nolint:errcheck

				switch k {
				case requestAuth:
					b, _ := json.Marshal(map[string]string{
						"access_token": "token",
					})
					f.Write(b) //nolint:errcheck
				case requestForms,
					requestForm,
					requestTask,
					requestCatalogs,
					requestCatalog,
					requestContacts,
					requestMembers,
					requestRoles,
					requestProfile,
					requestLists,
					requestListsTasks,
					requestInbox:
					if err := performRequest(accessToken, method, path, nil, nil, f); err != nil {
						log.Println(err)
						return
					}
				case requestRegistry:
					if err := performRequest(accessToken, method, path, nil, &RegistryRequest{}, f); err != nil {
						log.Println(err)
						return
					}
				}
			}()
		}
	}

	r := chi.NewRouter()
	for k, v := range requests {
		k := k
		mp := strings.Split(k, ":")
		method := mp[0]
		path := mp[1]

		v := v
		r.Method(method, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch k {
			case requestDownloadFile:
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Content-Disposition", "attachment; filename=\"index.html\"; filename*=UTF-8''index.html")
			case requestAddCallDetails, requestRegisterCallEvent:
				return
			default:
				w.Header().Set("Content-Type", "application/json")
			}

			token := r.Header.Get("Authorization")
			if k != requestAuth && token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				b, _ := json.Marshal(map[string]string{
					"error":      "Неверный токен авторизации.",
					"error_code": "invalid_token",
				})
				w.Write(b) //nolint:errcheck
				return
			} else if k != requestAuth && token != "Bearer token" {
				w.WriteHeader(http.StatusUnauthorized)
				b, _ := json.Marshal(map[string]string{
					"error":      "Токен авторизации не указан.",
					"error_code": "token_not_specified",
				})
				w.Write(b) //nolint:errcheck
				return
			}

			f, err := os.Open(v)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer f.Close() //nolint:errcheck

			if _, err := io.Copy(w, f); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}))
	}
	ts = httptest.NewServer(r)

	c, err := NewClient("login", "securityKey", WithBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	if err != nil {
		log.Fatalln(err)
	}

	cl = c
	os.Exit(m.Run())
}

func performRequest(accessToken, method, path string, q *url.Values, reqBody interface{}, respBody io.Writer) error {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		logger.Error("Error while parsing a URL!", zap.Error(err))
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
			logger.Error("Error while creating a new form file!", zap.Error(err))
			return err
		}
		if _, err := io.Copy(fw, reqBody.(*fileRequest).Reader); err != nil {
			logger.Error("Error while writing a file!", zap.Error(err))
			return err
		}
		if err := w.Close(); err != nil {
			logger.Error("Error while trying to close multipart writer!", zap.Error(err))
			return err
		}

		req, reqErr = http.NewRequest(method, u.String(), buf)
		contentTypeHeader = w.FormDataContentType()
	} else if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
			logger.Error("Error while encoding JSON!", zap.Error(err))
			return err
		}

		req, reqErr = http.NewRequest(method, u.String(), buf)
	} else {
		req, reqErr = http.NewRequest(method, u.String(), nil)
	}
	if reqErr != nil {
		logger.Error("Error while creating a request!", zap.Error(err))
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentTypeHeader)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Error while doing a request!", zap.Error(err))
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if _, err := io.Copy(respBody, resp.Body); err != nil {
		logger.Error("Error while copying a request!", zap.Error(err))
		return err
	}

	return nil
}

func TestNewClient(t *testing.T) {
	noopLogger := &noopLogger{}
	noopLogger.Error("test", errors.New("fake error"))

	zapLogger := &zapLogger{logger: zap.NewNop()}
	zapLogger.Error("test", errors.New("fake error"))

	testErr := Error{
		Code:        ErrCannotAddExternalUser,
		Description: "TEST",
	}
	assert.Equal(t, "API error: TEST (cannot_add_external_user)", testErr.Error())

	c, err := NewClient(
		pyrusLogin,
		pyrusSecurityKey,
		WithBaseURL(ts.URL),
		WithHTTPClient(ts.Client()),
		WithLogger(noopLogger),
		WithZapLogger(logger),
		WithEventBufferSize(100),
	)
	require.NoError(t, err)
	assert.NotNil(t, c)

	profile, err := c.Profile()
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestClient_ListenWebhook(t *testing.T) {
	handler, events := cl.WebhookHandler()

	go func() {
		require.NoError(t, http.ListenAndServe("127.0.0.1:30000", handler))
	}()

	f, err := os.Open("testdata/event.json")
	require.NoError(t, err)
	defer f.Close() //nolint:errcheck

	b, err := io.ReadAll(f)
	require.NoError(t, err)

	hasher := hmac.New(sha1.New, []byte("securityKey"))
	_, err = hasher.Write(b)
	require.NoError(t, err)
	hash := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:30000", bytes.NewBuffer(b))
	require.NoError(t, err)
	req.Header.Set("X-Pyrus-Sig", hash)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, <-events)
}

func TestClient_Auth(t *testing.T) {
	token, err := cl.Auth(pyrusLogin, pyrusSecurityKey)
	require.NoError(t, err)
	assert.Equal(t, "token", token)
}

func TestClient_Forms(t *testing.T) {
	forms, err := cl.Forms()
	require.NoError(t, err)
	assert.NotNil(t, forms)
}

func TestClient_Form(t *testing.T) {
	form, err := cl.Form(formID)
	require.NoError(t, err)
	assert.NotNil(t, form)
}

func TestClient_Registry(t *testing.T) {
	n := time.Now()
	tasks, err := cl.Registry(formID, &RegistryRequest{
		SimpleFormat:    true,
		IncludeArchived: true,
		CreatedBefore:   &n,
	})
	require.NoError(t, err)
	assert.NotNil(t, tasks)

	tasks, err = cl.Registry(formID, &RegistryRequest{
		SimpleFormat:    true,
		IncludeArchived: true,
		CreatedBefore:   &n,
		FieldFilters: map[int]string{
			1: "Payment",
			2: "15,20",
			3: "gt15",
			4: "lt2017-01-01",
			5: "gt15.5,lt21",
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, tasks)
}

func TestClient_Task(t *testing.T) {
	task, err := cl.Task(taskID)
	require.NoError(t, err)
	assert.NotNil(t, task)
}

func TestClient_Catalogs(t *testing.T) {
	catalogs, err := cl.Catalogs()
	require.NoError(t, err)
	assert.NotNil(t, catalogs)
}

func TestClient_Catalog(t *testing.T) {
	catalog, err := cl.Catalog(catalogID)
	require.NoError(t, err)
	assert.NotNil(t, catalog)
}

func TestClient_Contacts(t *testing.T) {
	contacts, err := cl.Contacts()
	require.NoError(t, err)
	assert.NotNil(t, contacts)
}

func TestClient_Members(t *testing.T) {
	members, err := cl.Members()
	require.NoError(t, err)
	assert.NotNil(t, members)
}

func TestClient_Roles(t *testing.T) {
	roles, err := cl.Roles()
	require.NoError(t, err)
	assert.NotNil(t, roles)
}

func TestClient_Profile(t *testing.T) {
	profile, err := cl.Profile()
	require.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestClient_Lists(t *testing.T) {
	lists, err := cl.Lists()
	require.NoError(t, err)
	assert.NotNil(t, lists)
}

func TestClient_TaskList(t *testing.T) {
	taskList, err := cl.TaskList(listID, 200, true)
	require.NoError(t, err)
	assert.NotNil(t, taskList)
}

func TestClient_Inbox(t *testing.T) {
	inbox, err := cl.Inbox(200)
	require.NoError(t, err)
	assert.NotNil(t, inbox)
}

func TestClient_CreateTask(t *testing.T) {
	task, err := cl.CreateTask(&TaskRequest{
		Text: "Пример",
		Subscribers: []*Person{
			{
				ID: 123456,
			},
		},
		Attachments: []*Attachment{
			{AttachmentID: fileID},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, task)
}

func TestClient_CommentTask(t *testing.T) {
	task, err := cl.CommentTask(taskID, &TaskCommentRequest{
		Subject: "Пример заголовка задачи",
		Text:    "Пример текста задачи",
		FieldUpdates: []*FormField{
			{
				ID:    1,
				Type:  FieldTypeText,
				Value: "example.org",
			},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, task)
}

func TestClient_UploadFile(t *testing.T) {
	f, err := os.Open("testdata/uploaded_file.json")
	require.NoError(t, err)
	defer f.Close() //nolint:errcheck

	file, err := cl.UploadFile("uploaded_file.json", f)
	require.NoError(t, err)
	assert.NotNil(t, file)
}

func TestClient_DownloadFile(t *testing.T) {
	file, err := cl.DownloadFile(fileID)
	require.NoError(t, err)
	assert.NotNil(t, file)
}

func TestClient_CreateCatalog(t *testing.T) {
	catalog, err := cl.CreateCatalog("BotTest", []string{"Имя", "Адрес"}, []*CatalogItem{
		{
			Values: []string{
				"Василий",
				"Островского 5",
			},
		},
		{
			Values: []string{
				"Иван",
				"Эсперанто 34",
			},
		},
	})
	assert.NoError(t, err)

	_ = catalog
}

func TestClient_SyncCatalog(t *testing.T) {
	syncCatalog, err := cl.SyncCatalog(catalogID, false, []string{"Имя", "Адрес"}, []*CatalogItem{
		{
			Values: []string{
				"Василий",
				"Островского 5",
			},
		},
		{
			Values: []string{
				"Иван",
				"Эсперанто 34",
			},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, syncCatalog)
}

func TestClient_CreateMember(t *testing.T) {
	member, err := cl.CreateMember(&MemberRequest{
		FirstName: "Савелий Игоревич",
		LastName:  "Красовский",
		Email:     "krasovskiisi@sovcombank.ru",
		Position:  "Аналитик",
		Phone:     "+79000000000",
	})
	require.NoError(t, err)
	assert.NotNil(t, member)
}

func TestClient_UpdateMember(t *testing.T) {
	member, err := cl.UpdateMember(memberID, &MemberRequest{
		FirstName: "Савелий Игоревич",
		LastName:  "Красовский",
		Email:     "krasovskiisi@sovcombank.ru",
		Position:  "Аналитик",
		Phone:     "+79000000000",
	})
	require.NoError(t, err)
	assert.NotNil(t, member)
}

func TestClient_BlockMember(t *testing.T) {
	member, err := cl.BlockMember(memberID)
	require.NoError(t, err)
	assert.NotNil(t, member)
}

func TestClient_CreateRole(t *testing.T) {
	role, err := cl.CreateRole("Боты ИБ", []int{529072})
	require.NoError(t, err)
	assert.NotNil(t, role)
}

func TestClient_UpdateRole(t *testing.T) {
	role, err := cl.UpdateRole(roleID, "Боты ИБ", []int{529072}, nil, false)
	require.NoError(t, err)
	assert.NotNil(t, role)
}

func TestClient_RegisterCall(t *testing.T) {
	call, err := cl.RegisterCall(&RegisterCallRequest{
		From:            "+79500000000",
		IntegrationGUID: "5d8dc3d6-27e7-4cd4-a057-2b4f4d74e0a5",
	})
	require.NoError(t, err)
	assert.NotNil(t, call)
}

func TestClient_AddCallDetails(t *testing.T) {
	err := cl.AddCallDetails("5d8dc3d6-27e7-4cd4-a057-2b4f4d74e0a5", &AddCallDetailsRequest{
		Rating: 5,
	})
	require.NoError(t, err)
}

func TestClient_RegisterCallEvent(t *testing.T) {
	err := cl.RegisterCallEvent("5d8dc3d6-27e7-4cd4-a057-2b4f4d74e0a5", CallEventTypeShow, "")
	require.NoError(t, err)
}
