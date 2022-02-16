package pyrus

// FieldType is a type of Form field.
type FieldType string

const (
	FieldTypeText         FieldType = "text"
	FieldTypeMoney        FieldType = "money"
	FieldTypeNumber       FieldType = "number"
	FieldTypeDate         FieldType = "date"
	FieldTypeTime         FieldType = "time"
	FieldTypeCheckmark    FieldType = "checkmark"
	FieldTypeDueDate      FieldType = "due_date"
	FieldTypeDueDateTime  FieldType = "due_date_time"
	FieldTypeEmail        FieldType = "email"
	FieldTypePhone        FieldType = "phone"
	FieldTypeFlag         FieldType = "flag"
	FieldTypeStep         FieldType = "step"
	FieldTypeStatus       FieldType = "status"
	FieldTypeCreationDate FieldType = "creation_date"
	FieldTypeNote         FieldType = "note"

	FieldTypeCatalog        FieldType = "catalog"
	FieldTypeFile           FieldType = "file"
	FieldTypePerson         FieldType = "person"
	FieldTypeAuthor         FieldType = "author"
	FieldTypeTable          FieldType = "table"
	FieldTypeMultipleChoice FieldType = "multiple_choice"
	FieldTypeTitle          FieldType = "title"
	FieldTypeFormLink       FieldType = "form_link"
	FieldTypeProject        FieldType = "project"
)

// PersonType is a type of Person.
type PersonType string

const (
	PersonTypeUser PersonType = "user"
	PersonTypeBot  PersonType = "bot"
	PersonTypeRole PersonType = "role"
)

// ChannelType is a type of Channel.
type ChannelType string

const (
	ChannelTypeEmail     ChannelType = "email"
	ChannelTypeTelegram  ChannelType = "telegram"
	ChannelTypeFacebook  ChannelType = "facebook"
	ChannelTypeVK        ChannelType = "vk"
	ChannelTypeViber     ChannelType = "viber"
	ChannelTypeMobileApp ChannelType = "mobile_app"
	ChannelTypeWebWidget ChannelType = "web_widget"
	ChannelTypeMoySklad  ChannelType = "moy_sklad"
	ChannelTypeZadarma   ChannelType = "zadarma"
	ChannelTypeAmoCRM    ChannelType = "amo_crm"
)

// ChoiceType is a type of approval choice in case of task.
type ChoiceType string

const (
	ChoiceTypeApproved     ChoiceType = "approved"
	ChoiceTypeAcknowledged ChoiceType = "acknowledged"
	ChoiceTypeRejected     ChoiceType = "rejected"
	ChoiceTypeRevoked      ChoiceType = "revoked"
	ChoiceTypeWaiting      ChoiceType = "waiting"
)

// ActionType is a type of action in case of task.
type ActionType string

const (
	ActionTypeFinished ActionType = "finished"
	ActionTypeReopened ActionType = "reopened"
)

// CheckmarkType is a type of checkmark. It could be only checked or unchecked.
type CheckmarkType string

const (
	CheckmarkTypeChecked   CheckmarkType = "checked"
	CheckmarkTypeUnchecked CheckmarkType = "unchecked"
)

// FlagType is a type of flag. While checkmark could be only checked or unchecked, flag could also has none state.
type FlagType string

const (
	FlagTypeNone      FlagType = "none"
	FlagTypeChecked   FlagType = "checked"
	FlagTypeUnchecked FlagType = "unchecked"
)

// StatusType is a type of status in case of task.
type StatusType string

const (
	StatusTypeOpen   StatusType = "open"
	StatusTypeClosed StatusType = "closed"
)

// CatalogHeaderType is a type of CatalogHeader
type CatalogHeaderType string

const (
	CatalogHeaderTypeText     CatalogHeaderType = "text"
	CatalogHeaderTypeWorkflow CatalogHeaderType = "workflow"
)

// DisconnectPartyType is a type of disconnect party. Only relevant for calls API.
type DisconnectPartyType string

const (
	DisconnectPartyTypeAgent  DisconnectPartyType = "agent"
	DisconnectPartyTypeClient DisconnectPartyType = "client"
	DisconnectPartyTypeError  DisconnectPartyType = "error"
	DisconnectPartyTypeOther  DisconnectPartyType = "other"
)

// CallStatusType is a type of call. Only relevant for calls API.
type CallStatusType string

const (
	CallStatusTypeAnswered CallStatusType = "answered"
	CallStatusTypeNoAnswer CallStatusType = "no answer"
	CallStatusTypeBusy     CallStatusType = "busy"
	CallStatusTypeError    CallStatusType = "error"
	CallStatusTypeOther    CallStatusType = "other"
)

// CallEventType is a type of call event. Only relevant for calls API.
type CallEventType string

const (
	CallEventTypeShow CallEventType = "show"
)
