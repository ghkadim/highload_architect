import openapi_client
import openapi_client_dialog
from openapi_client.api import default_api
from openapi_client_dialog.api import default_api
from openapi_client.model.user_register_post_request import UserRegisterPostRequest
from openapi_client.model.login_post_request import LoginPostRequest


class User:
    def __init__(self,
                 first_name="first_name",
                 second_name="second_name",
                 password="password",
                 user_id=None):
        self._api = openapi_client.api.default_api.DefaultApi()
        self._dialog_api = None
        self.first_name = first_name
        self.second_name = second_name
        self.password = password
        self.user_id = user_id
        if user_id is None:
            self.register()

    @staticmethod
    def search(first_name, second_name):
        users = openapi_client.api.default_api.DefaultApi().\
            user_search_get(first_name, second_name)
        resp = []
        for u in users:
            resp.append(User(
                first_name=u.first_name,
                second_name=u.second_name,
                password="password",
                user_id=u.id,
            ))
        return resp

    def register(self):
        res = self.api.user_register_post(
            user_register_post_request=UserRegisterPostRequest(
                first_name=self.first_name,
                second_name=self.second_name,
                password=self.password,
            ))
        self.user_id = res.user_id
        return res

    def login(self):
        res = self.api.login_post(
            login_post_request=LoginPostRequest(
                id=self.user_id,
                password=self.password,
            ))

        conf = openapi_client.Configuration.get_default_copy()
        conf.access_token = res.token
        self._api = openapi_client.api.default_api.DefaultApi(openapi_client.ApiClient(conf))

        conf = openapi_client_dialog.Configuration.get_default_copy()
        conf.access_token = res.token
        self._dialog_api = openapi_client_dialog.api.default_api.DefaultApi(openapi_client_dialog.ApiClient(conf))
        return res

    @property
    def api(self):
        return self._api

    @property
    def dialog_api(self):
        return self._dialog_api

    @property
    def access_token(self):
        return self._api.api_client.configuration.access_token

