# ModelsOrgCompat


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**avatar_url** | **str** |  | [optional] 
**created_at** | **str** |  | [optional] 
**deleted_at** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**settings** | **Dict[str, object]** |  | [optional] 
**slug** | **str** |  | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.models_org_compat import ModelsOrgCompat

# TODO update the JSON string below
json = "{}"
# create an instance of ModelsOrgCompat from a JSON string
models_org_compat_instance = ModelsOrgCompat.from_json(json)
# print the JSON string representation of the object
print(ModelsOrgCompat.to_json())

# convert the object into a dict
models_org_compat_dict = models_org_compat_instance.to_dict()
# create an instance of ModelsOrgCompat from a dict
models_org_compat_from_dict = ModelsOrgCompat.from_dict(models_org_compat_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


