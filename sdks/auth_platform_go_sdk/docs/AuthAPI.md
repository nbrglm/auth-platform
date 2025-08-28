# \AuthAPI

All URIs are relative to *http://localhost:3360*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiAuthFlowFlowIdGet**](AuthAPI.md#ApiAuthFlowFlowIdGet) | **Get** /api/auth/flow/{flowId} | Get User Flow Data
[**ApiAuthLoginPost**](AuthAPI.md#ApiAuthLoginPost) | **Post** /api/auth/login | User Login
[**ApiAuthLogoutPost**](AuthAPI.md#ApiAuthLogoutPost) | **Post** /api/auth/logout | Logout user
[**ApiAuthRefreshPost**](AuthAPI.md#ApiAuthRefreshPost) | **Post** /api/auth/refresh | Refresh Token
[**ApiAuthSignupPost**](AuthAPI.md#ApiAuthSignupPost) | **Post** /api/auth/signup | User Signup
[**ApiAuthVerifyEmailSendPost**](AuthAPI.md#ApiAuthVerifyEmailSendPost) | **Post** /api/auth/verify-email/send | Send Verification Email
[**ApiAuthVerifyEmailVerifyPost**](AuthAPI.md#ApiAuthVerifyEmailVerifyPost) | **Post** /api/auth/verify-email/verify | Verify Email Token



## ApiAuthFlowFlowIdGet

> HandlersUserFlowData ApiAuthFlowFlowIdGet(ctx, flowId).Execute()

Get User Flow Data



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	flowId := "flowId_example" // string | Flow ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthFlowFlowIdGet(context.Background(), flowId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthFlowFlowIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthFlowFlowIdGet`: HandlersUserFlowData
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthFlowFlowIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**flowId** | **string** | Flow ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthFlowFlowIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**HandlersUserFlowData**](HandlersUserFlowData.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthLoginPost

> HandlersUserLoginResult ApiAuthLoginPost(ctx).Data(data).Execute()

User Login



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	data := *openapiclient.NewHandlersUserLoginData("Email_example", "Password_example") // HandlersUserLoginData | User Login Data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthLoginPost(context.Background()).Data(data).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthLoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthLoginPost`: HandlersUserLoginResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthLoginPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthLoginPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersUserLoginData**](HandlersUserLoginData.md) | User Login Data | 

### Return type

[**HandlersUserLoginResult**](HandlersUserLoginResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthLogoutPost

> HandlersLogoutResult ApiAuthLogoutPost(ctx).XNAPSessionToken(xNAPSessionToken).XNAPRefreshToken(xNAPRefreshToken).Execute()

Logout user



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	xNAPSessionToken := "xNAPSessionToken_example" // string | Session token (optional)
	xNAPRefreshToken := "xNAPRefreshToken_example" // string | Refresh token (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthLogoutPost(context.Background()).XNAPSessionToken(xNAPSessionToken).XNAPRefreshToken(xNAPRefreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthLogoutPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthLogoutPost`: HandlersLogoutResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthLogoutPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthLogoutPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xNAPSessionToken** | **string** | Session token | 
 **xNAPRefreshToken** | **string** | Refresh token | 

### Return type

[**HandlersLogoutResult**](HandlersLogoutResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthRefreshPost

> HandlersRefreshTokenResult ApiAuthRefreshPost(ctx).XNAPRefreshToken(xNAPRefreshToken).Execute()

Refresh Token



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	xNAPRefreshToken := "xNAPRefreshToken_example" // string | Refresh token

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthRefreshPost(context.Background()).XNAPRefreshToken(xNAPRefreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthRefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthRefreshPost`: HandlersRefreshTokenResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthRefreshPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthRefreshPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xNAPRefreshToken** | **string** | Refresh token | 

### Return type

[**HandlersRefreshTokenResult**](HandlersRefreshTokenResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthSignupPost

> HandlersUserSignupResult ApiAuthSignupPost(ctx).Data(data).Execute()

User Signup



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	data := *openapiclient.NewHandlersUserSignupData("ConfirmPassword_example", "Email_example", "FirstName_example", "LastName_example", "Password_example") // HandlersUserSignupData | User Signup Data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthSignupPost(context.Background()).Data(data).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthSignupPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthSignupPost`: HandlersUserSignupResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthSignupPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthSignupPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersUserSignupData**](HandlersUserSignupData.md) | User Signup Data | 

### Return type

[**HandlersUserSignupResult**](HandlersUserSignupResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthVerifyEmailSendPost

> HandlersSendVerificationEmailResult ApiAuthVerifyEmailSendPost(ctx).Data(data).Execute()

Send Verification Email



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	data := *openapiclient.NewHandlersSendVerificationEmailData("Email_example") // HandlersSendVerificationEmailData | Send Verification Email Data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthVerifyEmailSendPost(context.Background()).Data(data).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthVerifyEmailSendPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthVerifyEmailSendPost`: HandlersSendVerificationEmailResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthVerifyEmailSendPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthVerifyEmailSendPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersSendVerificationEmailData**](HandlersSendVerificationEmailData.md) | Send Verification Email Data | 

### Return type

[**HandlersSendVerificationEmailResult**](HandlersSendVerificationEmailResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAuthVerifyEmailVerifyPost

> HandlersVerifyEmailTokenResult ApiAuthVerifyEmailVerifyPost(ctx).Data(data).Execute()

Verify Email Token



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	data := *openapiclient.NewHandlersVerifyEmailTokenData("Token_example") // HandlersVerifyEmailTokenData | Verify Email Token Data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ApiAuthVerifyEmailVerifyPost(context.Background()).Data(data).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ApiAuthVerifyEmailVerifyPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthVerifyEmailVerifyPost`: HandlersVerifyEmailTokenResult
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ApiAuthVerifyEmailVerifyPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthVerifyEmailVerifyPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersVerifyEmailTokenData**](HandlersVerifyEmailTokenData.md) | Verify Email Token Data | 

### Return type

[**HandlersVerifyEmailTokenResult**](HandlersVerifyEmailTokenResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

