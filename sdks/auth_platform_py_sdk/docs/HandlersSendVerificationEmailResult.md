# HandlersSendVerificationEmailResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**success** | **bool** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_send_verification_email_result import HandlersSendVerificationEmailResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersSendVerificationEmailResult from a JSON string
handlers_send_verification_email_result_instance = HandlersSendVerificationEmailResult.from_json(json)
# print the JSON string representation of the object
print(HandlersSendVerificationEmailResult.to_json())

# convert the object into a dict
handlers_send_verification_email_result_dict = handlers_send_verification_email_result_instance.to_dict()
# create an instance of HandlersSendVerificationEmailResult from a dict
handlers_send_verification_email_result_from_dict = HandlersSendVerificationEmailResult.from_dict(handlers_send_verification_email_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


