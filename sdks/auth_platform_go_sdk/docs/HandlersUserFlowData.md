# HandlersUserFlowData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**ExpiresAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**InviteToken** | Pointer to **string** | For Invite Flow | [optional] 
**MfaRequired** | Pointer to **bool** |  | [optional] 
**MfaVerified** | Pointer to **bool** |  | [optional] 
**Orgs** | Pointer to [**[]ModelsOrgCompat**](ModelsOrgCompat.md) |  | [optional] 
**ReturnTo** | Pointer to **string** | URL to redirect after flow completion | [optional] 
**SsoProvider** | Pointer to **string** | For SSO Flow, e.g., \&quot;google\&quot;, \&quot;github\&quot;, etc. | [optional] 
**SsoUserId** | Pointer to **string** | For SSO Flow, External User ID | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 

## Methods

### NewHandlersUserFlowData

`func NewHandlersUserFlowData() *HandlersUserFlowData`

NewHandlersUserFlowData instantiates a new HandlersUserFlowData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHandlersUserFlowDataWithDefaults

`func NewHandlersUserFlowDataWithDefaults() *HandlersUserFlowData`

NewHandlersUserFlowDataWithDefaults instantiates a new HandlersUserFlowData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *HandlersUserFlowData) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *HandlersUserFlowData) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *HandlersUserFlowData) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *HandlersUserFlowData) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetEmail

`func (o *HandlersUserFlowData) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *HandlersUserFlowData) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *HandlersUserFlowData) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *HandlersUserFlowData) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetExpiresAt

`func (o *HandlersUserFlowData) GetExpiresAt() string`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *HandlersUserFlowData) GetExpiresAtOk() (*string, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *HandlersUserFlowData) SetExpiresAt(v string)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *HandlersUserFlowData) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetId

`func (o *HandlersUserFlowData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *HandlersUserFlowData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *HandlersUserFlowData) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *HandlersUserFlowData) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInviteToken

`func (o *HandlersUserFlowData) GetInviteToken() string`

GetInviteToken returns the InviteToken field if non-nil, zero value otherwise.

### GetInviteTokenOk

`func (o *HandlersUserFlowData) GetInviteTokenOk() (*string, bool)`

GetInviteTokenOk returns a tuple with the InviteToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInviteToken

`func (o *HandlersUserFlowData) SetInviteToken(v string)`

SetInviteToken sets InviteToken field to given value.

### HasInviteToken

`func (o *HandlersUserFlowData) HasInviteToken() bool`

HasInviteToken returns a boolean if a field has been set.

### GetMfaRequired

`func (o *HandlersUserFlowData) GetMfaRequired() bool`

GetMfaRequired returns the MfaRequired field if non-nil, zero value otherwise.

### GetMfaRequiredOk

`func (o *HandlersUserFlowData) GetMfaRequiredOk() (*bool, bool)`

GetMfaRequiredOk returns a tuple with the MfaRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaRequired

`func (o *HandlersUserFlowData) SetMfaRequired(v bool)`

SetMfaRequired sets MfaRequired field to given value.

### HasMfaRequired

`func (o *HandlersUserFlowData) HasMfaRequired() bool`

HasMfaRequired returns a boolean if a field has been set.

### GetMfaVerified

`func (o *HandlersUserFlowData) GetMfaVerified() bool`

GetMfaVerified returns the MfaVerified field if non-nil, zero value otherwise.

### GetMfaVerifiedOk

`func (o *HandlersUserFlowData) GetMfaVerifiedOk() (*bool, bool)`

GetMfaVerifiedOk returns a tuple with the MfaVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaVerified

`func (o *HandlersUserFlowData) SetMfaVerified(v bool)`

SetMfaVerified sets MfaVerified field to given value.

### HasMfaVerified

`func (o *HandlersUserFlowData) HasMfaVerified() bool`

HasMfaVerified returns a boolean if a field has been set.

### GetOrgs

`func (o *HandlersUserFlowData) GetOrgs() []ModelsOrgCompat`

GetOrgs returns the Orgs field if non-nil, zero value otherwise.

### GetOrgsOk

`func (o *HandlersUserFlowData) GetOrgsOk() (*[]ModelsOrgCompat, bool)`

GetOrgsOk returns a tuple with the Orgs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgs

`func (o *HandlersUserFlowData) SetOrgs(v []ModelsOrgCompat)`

SetOrgs sets Orgs field to given value.

### HasOrgs

`func (o *HandlersUserFlowData) HasOrgs() bool`

HasOrgs returns a boolean if a field has been set.

### GetReturnTo

`func (o *HandlersUserFlowData) GetReturnTo() string`

GetReturnTo returns the ReturnTo field if non-nil, zero value otherwise.

### GetReturnToOk

`func (o *HandlersUserFlowData) GetReturnToOk() (*string, bool)`

GetReturnToOk returns a tuple with the ReturnTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReturnTo

`func (o *HandlersUserFlowData) SetReturnTo(v string)`

SetReturnTo sets ReturnTo field to given value.

### HasReturnTo

`func (o *HandlersUserFlowData) HasReturnTo() bool`

HasReturnTo returns a boolean if a field has been set.

### GetSsoProvider

`func (o *HandlersUserFlowData) GetSsoProvider() string`

GetSsoProvider returns the SsoProvider field if non-nil, zero value otherwise.

### GetSsoProviderOk

`func (o *HandlersUserFlowData) GetSsoProviderOk() (*string, bool)`

GetSsoProviderOk returns a tuple with the SsoProvider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSsoProvider

`func (o *HandlersUserFlowData) SetSsoProvider(v string)`

SetSsoProvider sets SsoProvider field to given value.

### HasSsoProvider

`func (o *HandlersUserFlowData) HasSsoProvider() bool`

HasSsoProvider returns a boolean if a field has been set.

### GetSsoUserId

`func (o *HandlersUserFlowData) GetSsoUserId() string`

GetSsoUserId returns the SsoUserId field if non-nil, zero value otherwise.

### GetSsoUserIdOk

`func (o *HandlersUserFlowData) GetSsoUserIdOk() (*string, bool)`

GetSsoUserIdOk returns a tuple with the SsoUserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSsoUserId

`func (o *HandlersUserFlowData) SetSsoUserId(v string)`

SetSsoUserId sets SsoUserId field to given value.

### HasSsoUserId

`func (o *HandlersUserFlowData) HasSsoUserId() bool`

HasSsoUserId returns a boolean if a field has been set.

### GetType

`func (o *HandlersUserFlowData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *HandlersUserFlowData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *HandlersUserFlowData) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *HandlersUserFlowData) HasType() bool`

HasType returns a boolean if a field has been set.

### GetUserId

`func (o *HandlersUserFlowData) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *HandlersUserFlowData) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *HandlersUserFlowData) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *HandlersUserFlowData) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


