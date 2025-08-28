# HandlersRefreshTokenResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tokens** | [**TokensTokens**](TokensTokens.md) |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_refresh_token_result import HandlersRefreshTokenResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersRefreshTokenResult from a JSON string
handlers_refresh_token_result_instance = HandlersRefreshTokenResult.from_json(json)
# print the JSON string representation of the object
print(HandlersRefreshTokenResult.to_json())

# convert the object into a dict
handlers_refresh_token_result_dict = handlers_refresh_token_result_instance.to_dict()
# create an instance of HandlersRefreshTokenResult from a dict
handlers_refresh_token_result_from_dict = HandlersRefreshTokenResult.from_dict(handlers_refresh_token_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


