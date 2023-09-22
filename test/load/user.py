import io
import random
import itertools

import openapi_client
import requests
from openapi_client.api import default_api
from openapi_client.model.user_register_post_request import UserRegisterPostRequest
from openapi_client.model.login_post_request import LoginPostRequest
from openapi_client.rest import RESTClientObject

from locust import HttpUser, events
from locust.clients import HttpSession

_user_num = itertools.count()
_user_ids = list()


class OpenApiUser(HttpUser):

    abstract = True
    """If abstract is True, the class is meant to be subclassed, and users will not choose this locust during a test"""

    def __init__(self,
                 *args,
                 first_name=None,
                 second_name=None,
                 password="password",
                 **kwargs):
        self.host = "http://localhost:8080"
        super(OpenApiUser, self).__init__(*args, **kwargs)
        self.openapi_client = openapi_client.api.default_api.DefaultApi()
        self.openapi_client.api_client.rest_client = Adapter(self.client)
        fn, sn = self._user_name(next(_user_num))
        self.first_name = fn if first_name is None else first_name
        self.second_name = sn if second_name is None else second_name
        self.password = password
        self.access_token = None
        self.user_id = None
        self.friends = list()
        self.register()

    @staticmethod
    def _user_name(num):
        number = str(num).zfill(5)
        return 'first_name_' + number, 'second_name_' + number

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
        _user_ids.append(self.user_id)
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

    def peek_friend(self, max_friend_num):
        if len(self.friends) >= max_friend_num:
            return random.choice(self.friends)
        while True:
            fr_id = random.choice(running_users())
            if fr_id == self.user_id:
                return None
            self.friends.append(fr_id)
            return fr_id


def running_users():
    return _user_ids


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


@events.init_command_line_parser.add_listener
def _(parser):
    parser.add_argument("--friends", type=int, env_var="FRIENDS",
                        dest="num_friends", default=2, help="Number of friends")


