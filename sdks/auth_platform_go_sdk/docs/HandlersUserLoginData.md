# HandlersUserLoginData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**FlowReturnTo** | Pointer to **string** | Optional field to store in the flow data which can be fetched by the client after login This can be used to redirect the user to a specific page after login or to maintain the state of the application. It is recommended to validate this field on the client side to prevent open redirect vulnerabilities. | [optional] 
**Password** | **string** |  | 

## Methods

### NewHandlersUserLoginData

`func NewHandlersUserLoginData(email string, password string, ) *HandlersUserLoginData`

NewHandlersUserLoginData instantiates a new HandlersUserLoginData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHandlersUserLoginDataWithDefaults

`func NewHandlersUserLoginDataWithDefaults() *HandlersUserLoginData`

NewHandlersUserLoginDataWithDefaults instantiates a new HandlersUserLoginData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *HandlersUserLoginData) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *HandlersUserLoginData) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *HandlersUserLoginData) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetFlowReturnTo

`func (o *HandlersUserLoginData) GetFlowReturnTo() string`

GetFlowReturnTo returns the FlowReturnTo field if non-nil, zero value otherwise.

### GetFlowReturnToOk

`func (o *HandlersUserLoginData) GetFlowReturnToOk() (*string, bool)`

GetFlowReturnToOk returns a tuple with the FlowReturnTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlowReturnTo

`func (o *HandlersUserLoginData) SetFlowReturnTo(v string)`

SetFlowReturnTo sets FlowReturnTo field to given value.

### HasFlowReturnTo

`func (o *HandlersUserLoginData) HasFlowReturnTo() bool`

HasFlowReturnTo returns a boolean if a field has been set.

### GetPassword

`func (o *HandlersUserLoginData) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *HandlersUserLoginData) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *HandlersUserLoginData) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


