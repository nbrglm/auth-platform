# ModelsErrorResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **int32** | Code is an error code that can be used for programmatic handling of errors. | [optional] 
**Debug** | Pointer to **string** | DebugMessage is a technical message that can be used for debugging. | [optional] 
**Message** | Pointer to **string** | Message is a user-friendly message that can be displayed to the end user. | [optional] 

## Methods

### NewModelsErrorResponse

`func NewModelsErrorResponse() *ModelsErrorResponse`

NewModelsErrorResponse instantiates a new ModelsErrorResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewModelsErrorResponseWithDefaults

`func NewModelsErrorResponseWithDefaults() *ModelsErrorResponse`

NewModelsErrorResponseWithDefaults instantiates a new ModelsErrorResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *ModelsErrorResponse) GetCode() int32`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ModelsErrorResponse) GetCodeOk() (*int32, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ModelsErrorResponse) SetCode(v int32)`

SetCode sets Code field to given value.

### HasCode

`func (o *ModelsErrorResponse) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetDebug

`func (o *ModelsErrorResponse) GetDebug() string`

GetDebug returns the Debug field if non-nil, zero value otherwise.

### GetDebugOk

`func (o *ModelsErrorResponse) GetDebugOk() (*string, bool)`

GetDebugOk returns a tuple with the Debug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebug

`func (o *ModelsErrorResponse) SetDebug(v string)`

SetDebug sets Debug field to given value.

### HasDebug

`func (o *ModelsErrorResponse) HasDebug() bool`

HasDebug returns a boolean if a field has been set.

### GetMessage

`func (o *ModelsErrorResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ModelsErrorResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ModelsErrorResponse) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ModelsErrorResponse) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


