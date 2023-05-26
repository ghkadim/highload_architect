import pytest
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from openapi_client.model.dialog_message import DialogMessage


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.dialog_user_id_list_get("12345")

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))


def test_unknown_user(default_user):
    with pytest.raises(Exception):
        default_user.api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get("12345")
    assert len(messages) == 0


def test_dialog(default_user, make_user):
    friend = make_user()

    default_user.api.dialog_user_id_send_post(
        friend.user_id,
        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get(friend.user_id)
    assert messages == [DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")]

    friend.login()

    friend.api.dialog_user_id_send_post(
        default_user.user_id,
        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get(friend.user_id)
    assert messages == [
        DialogMessage(_from=friend.user_id, to=default_user.user_id, text="hello"),
        DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")
    ]
