package pyrus

// ErrorCode is an "enum" for error codes.
// More about errors at:
// https://pyrus.com/en/help/api/errors-and-limits
type ErrorCode string

const (
	// ErrServerError is internal server error.
	// You are not supposed to see this error,
	// however if it persists for more than 10 minutes you may want to contact Pyrus team support.
	ErrServerError ErrorCode = "server_error"
	// ErrInvalidCredentials returns in case of invalid login or security key.
	ErrInvalidCredentials ErrorCode = "invalid_credentials"
	// ErrTokenNotSpecified returns if access_token was not specified in request.
	// Library sets access_token so you cannot get this error unless you are doing request by yourself.
	ErrTokenNotSpecified ErrorCode = "token_not_specified"
	// ErrRevokedToken returns if access_token has been revoked.
	ErrRevokedToken ErrorCode = "revoked_token"
	// ErrExpiredToken returns if access_token has expired.
	ErrExpiredToken ErrorCode = "expired_token"
	// ErrInvalidToken returns if access_token is invalid.
	ErrInvalidToken ErrorCode = "invalid_token"
	// ErrAuthorizationError returns in case of unknown authorization error.
	ErrAuthorizationError ErrorCode = "authorization_error"
	// ErrAccountBlocked returns if user account that executed the request is blocked.
	// You should contact your company’s administrator.
	ErrAccountBlocked ErrorCode = "account_blocked"
	// ErrInvalidFieldID returns if the field with the specified identifier does not exist in the form.
	// Please verify input parameters.
	ErrInvalidFieldID ErrorCode = "invalid_field_id"
	// ErrDeletedField returns if the field with the specified identifier has been deleted from the form.
	// Please verify input parameters.
	ErrDeletedField ErrorCode = "deleted_field"
	// ErrInvalidFieldName returns if the field with the specified name does not exist in the form.
	// Please verify input parameters.
	ErrInvalidFieldName ErrorCode = "invalid_field_name"
	// ErrInvalidFieldIDName returns if the field with the specified identifier and name does not exist or has been deleted.
	// Please verify input parameters.
	ErrInvalidFieldIDName ErrorCode = "invalid_field_id_name"
	// ErrNonUniqueName returns if the field name is not unique within the form.
	// Please use the form field identifier to write value into it.
	ErrNonUniqueName ErrorCode = "non_unique_name"
	// ErrFieldIdentityMissing returns if the form field identity (id or name) is not specified in the request.
	ErrFieldIdentityMissing ErrorCode = "field_identity_missing"
	// ErrDuplicateField returns in case of you are trying to modify the same field multiple times in one request.
	ErrDuplicateField ErrorCode = "duplicate_field"
	// ErrInvalidCatalogID returns if the catalog with the id specified in the form template does not exist.
	ErrInvalidCatalogID ErrorCode = "invalid_catalog_id"
	// ErrInvalidCatalogItemName returns if the item with the specified name does not exist in the catalog.
	ErrInvalidCatalogItemName ErrorCode = "invalid_catalog_item_name"
	// ErrNonUniqueCatalogItemName returns if there are multiple items with the specified name in the catalog.
	// Please use item identifier to set a value.
	ErrNonUniqueCatalogItemName ErrorCode = "non_unique_catalog_item_name"
	// ErrInvalidCatalogItemID returns if the item with the specified id does not exist in the catalog.
	ErrInvalidCatalogItemID ErrorCode = "invalid_catalog_item_id"
	// ErrCatalogItemIDNameMismatch returns if the item with the specified identifier doesn't have a specified value.
	ErrCatalogItemIDNameMismatch ErrorCode = "catalog_item_id_name_mismatch"
	// ErrInvalidEmail returns if the person with that specified email address does not exist.
	ErrInvalidEmail ErrorCode = "invalid_email"
	// ErrNonUniqueEmail returns if there are multiple persons with that specified email.
	ErrNonUniqueEmail ErrorCode = "non_unique_email"
	// ErrInvalidPersonID returns if the person with the specified id was not found.
	ErrInvalidPersonID ErrorCode = "invalid_person_id"
	// ErrInvalidPersonIDEmail returns if the person with the specified identifier has another email.
	ErrInvalidPersonIDEmail ErrorCode = "invalid_person_id_email"
	// ErrFormHasNoTask returns if there is no task with the specified id created based on the specified form template id.
	ErrFormHasNoTask ErrorCode = "form_has_no_task"
	// ErrUnrecognizedAttachmentID returns in case of invalid unique identifier of attachment.
	// User has no attachment with specified id.
	ErrUnrecognizedAttachmentID ErrorCode = "unrecognized_attachment_id"
	// ErrRequiredFieldMissing returns if one of the required form fields is missing.
	// The error description will indicate which one.
	ErrRequiredFieldMissing ErrorCode = "required_field_missing"
	// ErrTypeIsNotSupported returns if this field type does not support the writing of values.
	ErrTypeIsNotSupported ErrorCode = "type_is_not_supported"
	// ErrCatalogIdentityMissing means that catalog item_id must be specified in order to write value into the catalog field.
	ErrCatalogIdentityMissing ErrorCode = "catalog_identity_missing"
	// ErrIncorrectParametersCount points at incorrect parameter count for the selected filter operator.
	ErrIncorrectParametersCount ErrorCode = "incorrect_parameters_count"
	// ErrFilterTypeIsNotSupported returns if this field type is not supported as a filter value.
	ErrFilterTypeIsNotSupported ErrorCode = "filter_type_is_not_supported"
	// ErrStepFieldDoesNotExists return if there are no step fields in the form.
	// You can't filter this form by step number.
	ErrStepFieldDoesNotExists ErrorCode = "step_field_does_not_exists"
	// ErrCatalogItemIDMissing means that catalog item_id must be specified in order to write value.
	ErrCatalogItemIDMissing ErrorCode = "catalog_item_id_missing"
	// ErrPersonIdentityMissing means that person id or email must be specified in order to write value.
	ErrPersonIdentityMissing ErrorCode = "person_identity_missing"
	// ErrEitherDueDateOrDueCanBeSet returns if you set both due_date and due.
	ErrEitherDueDateOrDueCanBeSet ErrorCode = "either_due_date_or_due_can_be_set"
	// ErrNegativeDuration returns if you are trying to send negative duration.
	ErrNegativeDuration ErrorCode = "negative_duration"
	// ErrDurationIsTooLong returns if you are trying to send more than year duration.
	ErrDurationIsTooLong ErrorCode = "duration_is_too_long"
	// ErrDueMissing returns if duration was sent without due.
	ErrDueMissing ErrorCode = "due_missing"
	// ErrScheduledDateInPast returns if you are trying to schedule task in the past.
	ErrScheduledDateInPast ErrorCode = "scheduled_date_in_past"
	// ErrCannotAddFormProject returns if you are trying to attach a task to a form project or a form's subproject.
	ErrCannotAddFormProject ErrorCode = "cannot_add_form_project"
	// ErrFormTemplateCantBeRemovedFromTask returns because form template list can't be removed from the task.
	ErrFormTemplateCantBeRemovedFromTask ErrorCode = "form_template_cant_be_removed_from_task"
	// ErrNoFileInRequest returns if there are no files in the request.
	ErrNoFileInRequest ErrorCode = "no_file_in_request"
	// ErrTooLargeRequestLength returns if the file you are attaching exceeds the maximum allowable size (250MB).
	ErrTooLargeRequestLength ErrorCode = "too_large_request_length"
	// ErrRequiredParameterMissing returns if one of the required request parameters is missing.
	// The error description will indicate which one.
	ErrRequiredParameterMissing ErrorCode = "required_parameter_missing"
	// ErrTooManyTaskSteps returns if the maximum allowed number of task steps has been exceeded.
	ErrTooManyTaskSteps ErrorCode = "too_many_task_steps"
	// ErrInvalidValueFormat returns if the provided value can't be converted to the field type.
	// The error description will indicate the type and value.
	ErrInvalidValueFormat ErrorCode = "invalid_value_format"
	// ErrTooManyComments returns if the maximum allowed number of task comments exceeded (10000).
	ErrTooManyComments ErrorCode = "too_many_comments"
	// ErrInvalidStepNumber returns if you have passed a negative step number or zero.
	ErrInvalidStepNumber ErrorCode = "invalid_step_number"
	// ErrTaskLimitExceeded returns if the maximum allowed number of tasks for your organization exceeded.
	ErrTaskLimitExceeded ErrorCode = "task_limit_exceeded"
	// ErrFieldIsInTable returns if the field you are trying to change is a part of the table.
	// You can modify it only by modifying the table.
	ErrFieldIsInTable ErrorCode = "field_is_in_table"
	// ErrRequiredTableFieldMissing returns if the required field inside table is not filled.
	// The error text will contain the name of the table, the field name and the line number.
	ErrRequiredTableFieldMissing ErrorCode = "required_table_field_missing"
	// ErrDepartmentCatalogCanNotBeModified returns because you can not modify department catalog using public API.
	ErrDepartmentCatalogCanNotBeModified ErrorCode = "department_catalog_can_not_be_modified"
	// ErrCatalogDuplicateRows returns because catalog contains duplicate rows.
	// Remove them and try request again.
	ErrCatalogDuplicateRows ErrorCode = "catalog_duplicate_rows"
	// ErrEmptyCatalogHeaders returns because catalog headers can not be empty.
	ErrEmptyCatalogHeaders ErrorCode = "empty_catalog_headers"
	// ErrCanNotModifyDeletedCatalog returns because you are trying to update a catalog that was deleted.
	ErrCanNotModifyDeletedCatalog ErrorCode = "can_not_modify_deleted_catalog"
	// ErrCanNotModifyFirstColumn returns because you can not modify the first column in the catalog.
	ErrCanNotModifyFirstColumn ErrorCode = "can_not_modify_first_column"
	// ErrCatalogHeadersItemsMismatch returns because headers and values mismatch.
	ErrCatalogHeadersItemsMismatch ErrorCode = "catalog_headers_items_mismatch"
	// ErrTooManyCatalogItems returns if the maximum allowed catalog items count exceeded (15000).
	ErrTooManyCatalogItems ErrorCode = "too_many_catalog_items"
	// ErrCatalogItemMaxLengthExceeded returns if the maximum catalog item length exceeded (500).
	ErrCatalogItemMaxLengthExceeded ErrorCode = "catalog_item_max_length_exceeded"
	// ErrCatalogDuplicateHeaders returns because catalog contains duplicate headers.
	// Remove them and try request again.
	ErrCatalogDuplicateHeaders ErrorCode = "catalog_duplicate_headers"
	// ErrFormIDMissing returns if you are trying to pass empty form_id.
	ErrFormIDMissing ErrorCode = "form_id_missing"
	// ErrTextMissing returns if you are trying to pass empty text.
	ErrTextMissing ErrorCode = "text_missing"
	// ErrInvalidJSON returns if the request body is not a valid JSON.
	ErrInvalidJSON ErrorCode = "invalid_json"
	// ErrEmptyBody returns if the request body can't be empty.
	ErrEmptyBody ErrorCode = "empty_body"
	// ErrAccessDeniedProject returns because access to the requested project is denied.
	// Make sure that the user has all the required permissions.
	ErrAccessDeniedProject ErrorCode = "access_denied_project"
	// ErrAccessDeniedTask returns because access to the requested task is denied.
	// Make sure that the user has all the required permissions.
	ErrAccessDeniedTask ErrorCode = "access_denied_task"
	// ErrAccessDeniedCloseTask returns if you don't have enough permissions to close the task.
	ErrAccessDeniedCloseTask ErrorCode = "access_denied_close_task"
	// ErrAccessDeniedReopenTask returns if you don't have enough permissions to reopen the task.
	ErrAccessDeniedReopenTask ErrorCode = "access_denied_reopen_task"
	// ErrAccessDeniedCatalog returns because access to the requested catalog is denied.
	// Make sure that thw user has all the required permissions.
	ErrAccessDeniedCatalog ErrorCode = "access_denied_catalog"
	// ErrAccessDeniedForm returns because access to the requested form is denied.
	// Make sure that the user has all the required permissions.
	ErrAccessDeniedForm ErrorCode = "access_denied_form"
	// ErrAccessDeniedPerson returns if you can't collaborate with the specified person.
	// Make sure that you have this person in your contact list or send them an invitation.
	ErrAccessDeniedPerson ErrorCode = "access_denied_person"
	// ErrTooManyRequests returns if you have reached the limit of requests per 10 minutes.
	// Please wait and try again later.
	ErrTooManyRequests ErrorCode = "too_many_requests"
	// ErrEmptyFile returns if you are trying to upload empty files.
	ErrEmptyFile ErrorCode = "empty_file"
	// ErrBadMultipartContent returns if you are trying to send bad body.
	ErrBadMultipartContent ErrorCode = "bad_multipart_content"
	// ErrInvalidTableRow returns because you cannot reset rows that have been deleted, or have not been created.
	ErrInvalidTableRow ErrorCode = "invalid_table_row"
	// ErrCannotAddExternalUser returns because an employee from another organization cannot be added to the task.
	ErrCannotAddExternalUser ErrorCode = "cannot_add_external_user"
	// ErrUnrecognizedIntegrationGUID returns if an integration with this id does not exist.
	ErrUnrecognizedIntegrationGUID ErrorCode = "unrecognized_integration_guid"
	// ErrUnrecognizedCallGUID returns if a call with this id does not exist.
	ErrUnrecognizedCallGUID ErrorCode = "unrecognized_call_guid"
	// ErrUnsupportedAttachmentFormat returns if the attachment has unknown or unsupported audio extension.
	ErrUnsupportedAttachmentFormat ErrorCode = "unsupported_attachment_format"
)

// Error is a standard error returned by Pyrus API in case of any problem with request.
type Error struct {
	Code        ErrorCode `json:"error_code"`
	Description string    `json:"error"`

	// Returns in case of 404
	Message string `json:"Message"`
}

// Error returns error as a human readable string
func (e Error) Error() string {
	return "API error: " + e.Description + " (" + string(e.Code) + ")"
}
