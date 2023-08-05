import io

import openapi_client
import requests
from openapi_client.api import default_api
from openapi_client.model.user_register_post_request import UserRegisterPostRequest
from openapi_client.model.login_post_request import LoginPostRequest
from openapi_client.rest import RESTClientObject

from locust import HttpUser
from locust.clients import HttpSession


class OpenApiUser(HttpUser):

    abstract = True
    """If abstract is True, the class is meant to be subclassed, and users will not choose this locust during a test"""

    def __init__(self,
                 *args,
                 first_name="first_name",
                 second_name="second_name",
                 password="password",
                 **kwargs):
        self.host = "http://localhost:8080"
        super(OpenApiUser, self).__init__(*args, **kwargs)
        self.openapi_client = openapi_client.api.default_api.DefaultApi()
        self.openapi_client.api_client.rest_client = Adapter(self.client)
        self.first_name = first_name
        self.second_name = second_name
        self.password = password
        self.access_token = None
        self.register()
        # ids = self.find_user_id(self.first_name, self.second_name)
        # if len(ids) > 1:
        #     raise RuntimeError(f'more then one user with {self.first_name} {self.second_name}: {ids}')
        # elif len(ids) == 1:
        #     self.user_id = ids[0]
        # else:
        #     self.register()

    def find_user_id(self, first_name, second_name):
        with self.client.rename_request("/user/search"):
            users = self.openapi_client.user_search_get(first_name, second_name)
            return [u.id for u in users]

    def register(self):
        res = self.openapi_client.user_register_post(
            user_register_post_request=UserRegisterPostRequest(
                first_name=self.first_name,
                second_name=self.second_name,
                password=self.password,
            ))
        self.user_id = res.user_id
        return res

    def login(self):
        res = self.openapi_client.login_post(
            login_post_request=LoginPostRequest(
                id=self.user_id,
                password=self.password,
            ))

        conf = openapi_client.Configuration.get_default_copy()
        conf.access_token = res.token
        self.openapi_client = openapi_client.api.default_api.DefaultApi(openapi_client.ApiClient(conf))
        self.openapi_client.api_client.rest_client = Adapter(self.client)
        self.access_token = conf.access_token
        return res


class ResponseAdapter(io.IOBase):
    def __init__(self, resp: requests.Response):
        self.resp = resp
        self.status = resp.status_code
        self.reason = resp.reason
        self.data = resp.content

    def getheaders(self):
        return self.resp.headers

    def getheader(self, name, default=None):
        return self.resp.headers.get(name, default)


class Adapter(RESTClientObject):
    def __init__(self, sess: HttpSession):
        self.session = sess

    def request(self, method, url, query_params=None, headers=None,
                body=None, post_params=None, _preload_content=True,
                _request_timeout=None):

        return ResponseAdapter(
            self.session.request(method, url, params=query_params, headers=headers,
                                 data=post_params, json=body, timeout=_request_timeout))
