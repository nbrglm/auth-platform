# HandlersVerifyEmailTokenResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**success** | **bool** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_verify_email_token_result import HandlersVerifyEmailTokenResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersVerifyEmailTokenResult from a JSON string
handlers_verify_email_token_result_instance = HandlersVerifyEmailTokenResult.from_json(json)
# print the JSON string representation of the object
print(HandlersVerifyEmailTokenResult.to_json())

# convert the object into a dict
handlers_verify_email_token_result_dict = handlers_verify_email_token_result_instance.to_dict()
# create an instance of HandlersVerifyEmailTokenResult from a dict
handlers_verify_email_token_result_from_dict = HandlersVerifyEmailTokenResult.from_dict(handlers_verify_email_token_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


