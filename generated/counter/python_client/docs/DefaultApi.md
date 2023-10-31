# openapi_client_counter.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**counter_counter_id_get**](DefaultApi.md#counter_counter_id_get) | **GET** /counter/{counter_id} | 


# **counter_counter_id_get**
> CounterCounterIdGet200Response counter_counter_id_get(counter_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client_counter
from openapi_client_counter.api import default_api
from openapi_client_counter.model.counter_counter_id_get500_response import CounterCounterIdGet500Response
from openapi_client_counter.model.counter_counter_id_get200_response import CounterCounterIdGet200Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client_counter.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client_counter.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client_counter.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    counter_id = "counter_id_example" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.counter_counter_id_get(counter_id)
        pprint(api_response)
    except openapi_client_counter.ApiException as e:
        print("Exception when calling DefaultApi->counter_counter_id_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **counter_id** | **str**|  |

### Return type

[**CounterCounterIdGet200Response**](CounterCounterIdGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно отправлено сообщение |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

