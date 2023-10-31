import pytest
import openapi_client_counter
from openapi_client_counter.api import default_api
from openapi_client_dialog.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from user import User
from collections.abc import Callable


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client_counter.exceptions.UnauthorizedException):
        api.counter_counter_id_get("unread_messages")


def test_zero_counter(default_user: User):
    res = default_user.counter_api.counter_counter_id_get("unread_messages")
    assert res.value == 0


def test_one_dialog(default_user: User, make_user: Callable[[], User]):
    friend = make_user()
    friend.login()

    def check_counter(user_counter, friend_counter):
        res = default_user.counter_api.counter_counter_id_get("unread_messages")
        assert res.value == user_counter
        res = friend.counter_api.counter_counter_id_get("unread_messages")
        assert res.value == friend_counter

    msg1 = default_user.dialog_api.dialog_user_id_send_post(
        friend.user_id,
        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))
    check_counter(1, 1)

    default_user.dialog_api.dialog_message_message_id_read_put(msg1.message_id)
    check_counter(0, 1)


def test_multiple_dialogs(default_user: User, make_user: Callable[[], User]):
    friends = []
    for i in range(10):
        f = make_user()
        f.login()
        friends.append(f)

    messages = []
    for f in friends:
        msg = default_user.dialog_api.dialog_user_id_send_post(
            f.user_id,
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))
        messages.append(msg)

    res = default_user.counter_api.counter_counter_id_get("unread_messages")
    assert res.value == len(messages)

    for f in friends:
        res = f.counter_api.counter_counter_id_get("unread_messages")
        assert res.value == 1

    for i, msg in enumerate(messages):
        default_user.dialog_api.dialog_message_message_id_read_put(msg.message_id)
        res = default_user.counter_api.counter_counter_id_get("unread_messages")
        assert res.value == len(messages) - i - 1

    for f in friends:
        res = f.counter_api.counter_counter_id_get("unread_messages")
        assert res.value == 1
