import json

from contextlib import closing
from websocket import create_connection
from websocket._exceptions import WebSocketTimeoutException, WebSocketBadStatusException


class WebSocketException(Exception):
    def __init__(self, status_code):
        self.status_code = status_code


class Post:
    def __init__(self, postId=None, postText=None, author_user_id=None):
        self.id = postId
        self.text = postText
        self.author_id = author_user_id

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, self.__class__):
            return False

        return self.id == other.id and self.text == other.text and self.author_id == other.author_id

    def __repr__(self):
        return f'{self.id, self.text, self.author_id}'


class AsyncApi:
    def __init__(self, host):
        self.host = host

    def post_feed_posted(self, token=None, timeout=5):
        header = None
        if token is not None:
            header = {"Authorization": "Bearer " + token}

        url = f'{self.host}/post/feed/posted'

        class _Conn:
            def __init__(self):
                try:
                    self.conn = create_connection(url, timeout=timeout, header=header)
                except WebSocketBadStatusException as e:
                    raise WebSocketException(status_code=e.status_code)

            def __iter__(self):
                return self

            def __next__(self):
                try:
                    return Post(**json.loads(self.conn.recv()))
                except WebSocketBadStatusException as e:
                    raise WebSocketException(status_code=e.status_code)
                except WebSocketTimeoutException:
                    raise StopIteration

            def close(self):
                self.conn.close()

        return closing(_Conn())

