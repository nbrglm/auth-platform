# auth_platform_py_sdk.AuthApi

All URIs are relative to *http://localhost:3360*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_auth_flow_flow_id_get**](AuthApi.md#api_auth_flow_flow_id_get) | **GET** /api/auth/flow/{flowId} | Get User Flow Data
[**api_auth_login_post**](AuthApi.md#api_auth_login_post) | **POST** /api/auth/login | User Login
[**api_auth_logout_post**](AuthApi.md#api_auth_logout_post) | **POST** /api/auth/logout | Logout user
[**api_auth_refresh_post**](AuthApi.md#api_auth_refresh_post) | **POST** /api/auth/refresh | Refresh Token
[**api_auth_signup_post**](AuthApi.md#api_auth_signup_post) | **POST** /api/auth/signup | User Signup
[**api_auth_verify_email_send_post**](AuthApi.md#api_auth_verify_email_send_post) | **POST** /api/auth/verify-email/send | Send Verification Email
[**api_auth_verify_email_verify_post**](AuthApi.md#api_auth_verify_email_verify_post) | **POST** /api/auth/verify-email/verify | Verify Email Token


# **api_auth_flow_flow_id_get**
> HandlersUserFlowData api_auth_flow_flow_id_get(flow_id)

Get User Flow Data

Retrieves user flow data based on the provided flow ID.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_user_flow_data import HandlersUserFlowData
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    flow_id = 'flow_id_example' # str | Flow ID

    try:
        # Get User Flow Data
        api_response = api_instance.api_auth_flow_flow_id_get(flow_id)
        print("The response of AuthApi->api_auth_flow_flow_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_flow_flow_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **flow_id** | **str**| Flow ID | 

### Return type

[**HandlersUserFlowData**](HandlersUserFlowData.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | User Flow Data |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_login_post**
> HandlersUserLoginResult api_auth_login_post(data)

User Login

Handles user login requests.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_user_login_data import HandlersUserLoginData
from auth_platform_py_sdk.models.handlers_user_login_result import HandlersUserLoginResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    data = auth_platform_py_sdk.HandlersUserLoginData() # HandlersUserLoginData | User Login Data

    try:
        # User Login
        api_response = api_instance.api_auth_login_post(data)
        print("The response of AuthApi->api_auth_login_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_login_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersUserLoginData**](HandlersUserLoginData.md)| User Login Data | 

### Return type

[**HandlersUserLoginResult**](HandlersUserLoginResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | User Login Result |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**403** | Forbidden - Email Not Verified OR User does not belong to any organization |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_logout_post**
> HandlersLogoutResult api_auth_logout_post(x_nap_session_token=x_nap_session_token, x_nap_refresh_token=x_nap_refresh_token)

Logout user

Logs out the user by revoking their session using session token or refresh token. Requires atleast one of the tokens.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_logout_result import HandlersLogoutResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    x_nap_session_token = 'x_nap_session_token_example' # str | Session token (optional)
    x_nap_refresh_token = 'x_nap_refresh_token_example' # str | Refresh token (optional)

    try:
        # Logout user
        api_response = api_instance.api_auth_logout_post(x_nap_session_token=x_nap_session_token, x_nap_refresh_token=x_nap_refresh_token)
        print("The response of AuthApi->api_auth_logout_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_logout_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **x_nap_session_token** | **str**| Session token | [optional] 
 **x_nap_refresh_token** | **str**| Refresh token | [optional] 

### Return type

[**HandlersLogoutResult**](HandlersLogoutResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Logout result |  -  |
**400** | Bad Request - Invalid or missing tokens |  -  |
**401** | Unauthorized - Invalid or expired tokens |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_refresh_post**
> HandlersRefreshTokenResult api_auth_refresh_post(x_nap_refresh_token)

Refresh Token

Handles token refresh requests.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_refresh_token_result import HandlersRefreshTokenResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    x_nap_refresh_token = 'x_nap_refresh_token_example' # str | Refresh token

    try:
        # Refresh Token
        api_response = api_instance.api_auth_refresh_post(x_nap_refresh_token)
        print("The response of AuthApi->api_auth_refresh_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_refresh_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **x_nap_refresh_token** | **str**| Refresh token | 

### Return type

[**HandlersRefreshTokenResult**](HandlersRefreshTokenResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | New tokens |  -  |
**400** | Bad Request - Invalid or missing tokens |  -  |
**401** | Unauthorized - Invalid or expired tokens - Proceed to Login |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_signup_post**
> HandlersUserSignupResult api_auth_signup_post(data)

User Signup

Handles user registration requests.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_user_signup_data import HandlersUserSignupData
from auth_platform_py_sdk.models.handlers_user_signup_result import HandlersUserSignupResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    data = auth_platform_py_sdk.HandlersUserSignupData() # HandlersUserSignupData | User Signup Data

    try:
        # User Signup
        api_response = api_instance.api_auth_signup_post(data)
        print("The response of AuthApi->api_auth_signup_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_signup_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersUserSignupData**](HandlersUserSignupData.md)| User Signup Data | 

### Return type

[**HandlersUserSignupResult**](HandlersUserSignupResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | User Signup Result |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized - Invalid Invite Token or Missing Invite Token or Domain Not Allowed |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_verify_email_send_post**
> HandlersSendVerificationEmailResult api_auth_verify_email_send_post(data)

Send Verification Email

Sends a verification email to the user.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_send_verification_email_data import HandlersSendVerificationEmailData
from auth_platform_py_sdk.models.handlers_send_verification_email_result import HandlersSendVerificationEmailResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    data = auth_platform_py_sdk.HandlersSendVerificationEmailData() # HandlersSendVerificationEmailData | Send Verification Email Data

    try:
        # Send Verification Email
        api_response = api_instance.api_auth_verify_email_send_post(data)
        print("The response of AuthApi->api_auth_verify_email_send_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_verify_email_send_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersSendVerificationEmailData**](HandlersSendVerificationEmailData.md)| Send Verification Email Data | 

### Return type

[**HandlersSendVerificationEmailResult**](HandlersSendVerificationEmailResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Send Verification Email Result |  -  |
**400** | Bad Request - Invalid Input or User does not exist |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_auth_verify_email_verify_post**
> HandlersVerifyEmailTokenResult api_auth_verify_email_verify_post(data)

Verify Email Token

Verifies the email using the provided token.

### Example


```python
import auth_platform_py_sdk
from auth_platform_py_sdk.models.handlers_verify_email_token_data import HandlersVerifyEmailTokenData
from auth_platform_py_sdk.models.handlers_verify_email_token_result import HandlersVerifyEmailTokenResult
from auth_platform_py_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:3360
# See configuration.py for a list of all supported configuration parameters.
configuration = auth_platform_py_sdk.Configuration(
    host = "http://localhost:3360"
)


# Enter a context with an instance of the API client
with auth_platform_py_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = auth_platform_py_sdk.AuthApi(api_client)
    data = auth_platform_py_sdk.HandlersVerifyEmailTokenData() # HandlersVerifyEmailTokenData | Verify Email Token Data

    try:
        # Verify Email Token
        api_response = api_instance.api_auth_verify_email_verify_post(data)
        print("The response of AuthApi->api_auth_verify_email_verify_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->api_auth_verify_email_verify_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **data** | [**HandlersVerifyEmailTokenData**](HandlersVerifyEmailTokenData.md)| Verify Email Token Data | 

### Return type

[**HandlersVerifyEmailTokenResult**](HandlersVerifyEmailTokenResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Verify Email Token Result |  -  |
**400** | Bad Request - Invalid Input or Token |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

