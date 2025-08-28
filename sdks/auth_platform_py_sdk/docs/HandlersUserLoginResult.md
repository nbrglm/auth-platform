# HandlersUserLoginResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**flow_id** | **str** |  | [optional] 
**message** | **str** |  | [optional] 
**require_email_verification** | **bool** |  | [optional] 
**tokens** | [**TokensTokens**](TokensTokens.md) |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_user_login_result import HandlersUserLoginResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersUserLoginResult from a JSON string
handlers_user_login_result_instance = HandlersUserLoginResult.from_json(json)
# print the JSON string representation of the object
print(HandlersUserLoginResult.to_json())

# convert the object into a dict
handlers_user_login_result_dict = handlers_user_login_result_instance.to_dict()
# create an instance of HandlersUserLoginResult from a dict
handlers_user_login_result_from_dict = HandlersUserLoginResult.from_dict(handlers_user_login_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


