# TokensTokens

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefreshToken** | Pointer to **string** | RefreshToken is the generated refresh token.  This is base64.RawURLEncoding encoded. DO NOT DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT. | [optional] 
**RefreshTokenExpiry** | Pointer to **string** | RefreshTokenExpiry is the expiration time of the refresh token. | [optional] 
**SessionId** | Pointer to **string** | SessionId is the unique identifier for the session. | [optional] 
**SessionToken** | Pointer to **string** | SessionToken is the generated session token.  This is a jwt which is base64.RawURLEncoding encoded. YOU NEED TO DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT. DO NOT USE IT AS IS. VALIDATION WILL FAIL WITHOUT DECODING. ONLY WHEN DECODED, YOU SHOULD PASS IT TO THE THINGS THAT REQUIRE THE SESSION TOKEN. | [optional] 
**SessionTokenExpiry** | Pointer to **string** | SessionTokenExpiry is the expiration time of the session token. | [optional] 

## Methods

### NewTokensTokens

`func NewTokensTokens() *TokensTokens`

NewTokensTokens instantiates a new TokensTokens object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTokensTokensWithDefaults

`func NewTokensTokensWithDefaults() *TokensTokens`

NewTokensTokensWithDefaults instantiates a new TokensTokens object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRefreshToken

`func (o *TokensTokens) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *TokensTokens) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *TokensTokens) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *TokensTokens) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### GetRefreshTokenExpiry

`func (o *TokensTokens) GetRefreshTokenExpiry() string`

GetRefreshTokenExpiry returns the RefreshTokenExpiry field if non-nil, zero value otherwise.

### GetRefreshTokenExpiryOk

`func (o *TokensTokens) GetRefreshTokenExpiryOk() (*string, bool)`

GetRefreshTokenExpiryOk returns a tuple with the RefreshTokenExpiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenExpiry

`func (o *TokensTokens) SetRefreshTokenExpiry(v string)`

SetRefreshTokenExpiry sets RefreshTokenExpiry field to given value.

### HasRefreshTokenExpiry

`func (o *TokensTokens) HasRefreshTokenExpiry() bool`

HasRefreshTokenExpiry returns a boolean if a field has been set.

### GetSessionId

`func (o *TokensTokens) GetSessionId() string`

GetSessionId returns the SessionId field if non-nil, zero value otherwise.

### GetSessionIdOk

`func (o *TokensTokens) GetSessionIdOk() (*string, bool)`

GetSessionIdOk returns a tuple with the SessionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionId

`func (o *TokensTokens) SetSessionId(v string)`

SetSessionId sets SessionId field to given value.

### HasSessionId

`func (o *TokensTokens) HasSessionId() bool`

HasSessionId returns a boolean if a field has been set.

### GetSessionToken

`func (o *TokensTokens) GetSessionToken() string`

GetSessionToken returns the SessionToken field if non-nil, zero value otherwise.

### GetSessionTokenOk

`func (o *TokensTokens) GetSessionTokenOk() (*string, bool)`

GetSessionTokenOk returns a tuple with the SessionToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionToken

`func (o *TokensTokens) SetSessionToken(v string)`

SetSessionToken sets SessionToken field to given value.

### HasSessionToken

`func (o *TokensTokens) HasSessionToken() bool`

HasSessionToken returns a boolean if a field has been set.

### GetSessionTokenExpiry

`func (o *TokensTokens) GetSessionTokenExpiry() string`

GetSessionTokenExpiry returns the SessionTokenExpiry field if non-nil, zero value otherwise.

### GetSessionTokenExpiryOk

`func (o *TokensTokens) GetSessionTokenExpiryOk() (*string, bool)`

GetSessionTokenExpiryOk returns a tuple with the SessionTokenExpiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionTokenExpiry

`func (o *TokensTokens) SetSessionTokenExpiry(v string)`

SetSessionTokenExpiry sets SessionTokenExpiry field to given value.

### HasSessionTokenExpiry

`func (o *TokensTokens) HasSessionTokenExpiry() bool`

HasSessionTokenExpiry returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


