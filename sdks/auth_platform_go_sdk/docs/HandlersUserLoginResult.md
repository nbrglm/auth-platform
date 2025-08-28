# HandlersUserLoginResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FlowId** | Pointer to **string** |  | [optional] 
**Message** | Pointer to **string** |  | [optional] 
**RequireEmailVerification** | Pointer to **bool** |  | [optional] 
**Tokens** | Pointer to [**TokensTokens**](TokensTokens.md) |  | [optional] 

## Methods

### NewHandlersUserLoginResult

`func NewHandlersUserLoginResult() *HandlersUserLoginResult`

NewHandlersUserLoginResult instantiates a new HandlersUserLoginResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHandlersUserLoginResultWithDefaults

`func NewHandlersUserLoginResultWithDefaults() *HandlersUserLoginResult`

NewHandlersUserLoginResultWithDefaults instantiates a new HandlersUserLoginResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFlowId

`func (o *HandlersUserLoginResult) GetFlowId() string`

GetFlowId returns the FlowId field if non-nil, zero value otherwise.

### GetFlowIdOk

`func (o *HandlersUserLoginResult) GetFlowIdOk() (*string, bool)`

GetFlowIdOk returns a tuple with the FlowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlowId

`func (o *HandlersUserLoginResult) SetFlowId(v string)`

SetFlowId sets FlowId field to given value.

### HasFlowId

`func (o *HandlersUserLoginResult) HasFlowId() bool`

HasFlowId returns a boolean if a field has been set.

### GetMessage

`func (o *HandlersUserLoginResult) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *HandlersUserLoginResult) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *HandlersUserLoginResult) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *HandlersUserLoginResult) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetRequireEmailVerification

`func (o *HandlersUserLoginResult) GetRequireEmailVerification() bool`

GetRequireEmailVerification returns the RequireEmailVerification field if non-nil, zero value otherwise.

### GetRequireEmailVerificationOk

`func (o *HandlersUserLoginResult) GetRequireEmailVerificationOk() (*bool, bool)`

GetRequireEmailVerificationOk returns a tuple with the RequireEmailVerification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequireEmailVerification

`func (o *HandlersUserLoginResult) SetRequireEmailVerification(v bool)`

SetRequireEmailVerification sets RequireEmailVerification field to given value.

### HasRequireEmailVerification

`func (o *HandlersUserLoginResult) HasRequireEmailVerification() bool`

HasRequireEmailVerification returns a boolean if a field has been set.

### GetTokens

`func (o *HandlersUserLoginResult) GetTokens() TokensTokens`

GetTokens returns the Tokens field if non-nil, zero value otherwise.

### GetTokensOk

`func (o *HandlersUserLoginResult) GetTokensOk() (*TokensTokens, bool)`

GetTokensOk returns a tuple with the Tokens field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokens

`func (o *HandlersUserLoginResult) SetTokens(v TokensTokens)`

SetTokens sets Tokens field to given value.

### HasTokens

`func (o *HandlersUserLoginResult) HasTokens() bool`

HasTokens returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


