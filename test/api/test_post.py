import logging
import time

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


@pytest.mark.parametrize(
    "user_num,posts_per_user", [
        (2, 1),
        (4, 10),
        (11, 100)
    ])
def test_feed(default_user: User, make_user, user_num, posts_per_user):
    users = []
    for i in range(user_num):
        user = make_user()
        user.login()
        users.append(user)
        default_user.api.friend_set_user_id_put(user.user_id)

    posts = []
    post_idx = 0
    for i in range(posts_per_user):
        for user in users:
            text = f"post {post_idx} from {user.user_id}"
            post_id = user.api.post_create_post(
                post_create_post_request=PostCreatePostRequest(text=text))
            posts = [Post(id=post_id, text=text, author_user_id=user.user_id)] + posts
            post_idx += 1

    def cmp_feed(offset, limit, feed):
        retry_count = 3
        for _i in range(1, retry_count + 1):
            try:
                post_feed = default_user.api.post_feed_get(offset=offset, limit=limit)
                assert post_feed == feed[offset:limit+offset]
                return
            except AssertionError as e:
                if _i == retry_count:
                    raise e
                time.sleep(2)
                continue

    def cmp_feeds(feed):
        cmp_feed(0, 10, feed)
        cmp_feed(0, 40, feed)
        cmp_feed(0, 1000, feed)
        cmp_feed(0, 1010, feed)
        cmp_feed(10, 20, feed)
        cmp_feed(10, 1000, feed)

    cmp_feeds(posts)

    default_user.api.friend_delete_user_id_put(users[0].user_id)
    updated_posts = list(filter(lambda p: p.author_user_id != users[0].user_id, posts))
    cmp_feeds(updated_posts)

    default_user.api.friend_set_user_id_put(users[0].user_id)
    cmp_feeds(posts)

    for i in range(posts_per_user):
        for user in users:
            text = f"post {post_idx} from {user.user_id}"
            post_id = user.api.post_create_post(
                post_create_post_request=PostCreatePostRequest(text=text))
            posts = [Post(id=post_id, text=text, author_user_id=user.user_id)] + posts
            post_idx += 1

    cmp_feeds(posts)
