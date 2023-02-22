import pytest
import openapi_client
from user import User
from openapi_client.api import default_api
from openapi_client.model.post_create_post_request import PostCreatePostRequest
from openapi_client.model.post_update_put_request import PostUpdatePutRequest
from openapi_client.model.post import Post


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text="foo"))

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.post_update_put(
            post_update_put_request=PostUpdatePutRequest(id="123", text="bar"))

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.post_delete_id_put(id="123")

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.post_feed_get()


def test_post(default_user: User):
    post_id = default_user.api.post_create_post(
        post_create_post_request=PostCreatePostRequest(text="hello"))
    default_user.api.post_update_put(
        post_update_put_request=PostUpdatePutRequest(id=post_id, text="hello world"))

    post = default_user.api.post_get_id_get(post_id)
    assert post.id == post_id
    assert post.text == "hello world"
    assert post.author_user_id == default_user.user_id

    default_user.api.post_delete_id_put(id=post.id)

    with pytest.raises(openapi_client.exceptions.NotFoundException):
        default_user.api.post_get_id_get(post_id)


def test_empty_feed(default_user: User):
    assert default_user.api.post_feed_get(offset=0, limit=100) == []


def test_feed(default_user: User, make_user):
    users = []
    for i in range(0, 10):
        user = make_user()
        user.login()
        users.append(user)
        default_user.api.friend_set_user_id_put(user.user_id)

    posts = []
    for i in range(0, 2):
        for user in users:
            text = f"post {i} from {user.user_id}"
            post_id = user.api.post_create_post(
                post_create_post_request=PostCreatePostRequest(text=text))
            posts = [Post(id=post_id, text=text, author_user_id=user.user_id)] + posts

    post_feed = default_user.api.post_feed_get(offset=0, limit=40)
    assert post_feed == posts[:40]
    post_feed = default_user.api.post_feed_get(offset=0, limit=10)
    assert post_feed == posts[:10]
    post_feed = default_user.api.post_feed_get(offset=10, limit=10)
    assert post_feed == posts[10:20]

    default_user.api.friend_delete_user_id_put(users[0].user_id)
    updated_posts = list(filter(lambda p: p.author_user_id != users[0].user_id, posts))
    post_feed = default_user.api.post_feed_get(offset=0, limit=40)
    assert post_feed == updated_posts[:40]
    post_feed = default_user.api.post_feed_get(offset=0, limit=10)
    assert post_feed == updated_posts[:10]
    post_feed = default_user.api.post_feed_get(offset=10, limit=10)
    assert post_feed == updated_posts[10:20]

    default_user.api.friend_set_user_id_put(users[0].user_id)
    post_feed = default_user.api.post_feed_get(offset=0, limit=40)
    assert post_feed == posts[:40]
    post_feed = default_user.api.post_feed_get(offset=0, limit=10)
    assert post_feed == posts[:10]
    post_feed = default_user.api.post_feed_get(offset=10, limit=10)
    assert post_feed == posts[10:20]

    for i in range(0, 2):
        for user in users:
            text = f"post {i} from {user.user_id}"
            post_id = user.api.post_create_post(
                post_create_post_request=PostCreatePostRequest(text=text))
            posts = [Post(id=post_id, text=text, author_user_id=user.user_id)] + posts

    post_feed = default_user.api.post_feed_get(offset=0, limit=40)
    assert post_feed == posts[:40]
    post_feed = default_user.api.post_feed_get(offset=0, limit=10)
    assert post_feed == posts[:10]
    post_feed = default_user.api.post_feed_get(offset=10, limit=10)
    assert post_feed == posts[10:20]
