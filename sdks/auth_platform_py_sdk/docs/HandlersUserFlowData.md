# HandlersUserFlowData


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**email** | **str** |  | [optional] 
**expires_at** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**invite_token** | **str** | For Invite Flow | [optional] 
**mfa_required** | **bool** |  | [optional] 
**mfa_verified** | **bool** |  | [optional] 
**orgs** | [**List[ModelsOrgCompat]**](ModelsOrgCompat.md) |  | [optional] 
**return_to** | **str** | URL to redirect after flow completion | [optional] 
**sso_provider** | **str** | For SSO Flow, e.g., \&quot;google\&quot;, \&quot;github\&quot;, etc. | [optional] 
**sso_user_id** | **str** | For SSO Flow, External User ID | [optional] 
**type** | **str** |  | [optional] 
**user_id** | **str** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_user_flow_data import HandlersUserFlowData

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersUserFlowData from a JSON string
handlers_user_flow_data_instance = HandlersUserFlowData.from_json(json)
# print the JSON string representation of the object
print(HandlersUserFlowData.to_json())

# convert the object into a dict
handlers_user_flow_data_dict = handlers_user_flow_data_instance.to_dict()
# create an instance of HandlersUserFlowData from a dict
handlers_user_flow_data_from_dict = HandlersUserFlowData.from_dict(handlers_user_flow_data_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


