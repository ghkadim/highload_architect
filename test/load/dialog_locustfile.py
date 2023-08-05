import itertools
import random

from locust import events, run_single_user, constant_pacing, task, LoadTestShape
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from user import OpenApiUser

user_num = itertools.count()
user_ids = list()


class DialogUser(OpenApiUser):
    def __init__(self, env):
        self.user_num = next(user_num)
        fn, sn = DialogUser.user_name(self.user_num)
        super(DialogUser, self).__init__(env, first_name=fn, second_name=sn)
        self.login()
        self.num_dialog_friends = env.parsed_options.num_dialog_friends
        self.friends = list()
        user_ids.append(self.user_id)

    wait_time = constant_pacing(5)

    @staticmethod
    def user_name(num):
        number = str(num).zfill(5)
        return 'first_name_' + number, 'second_name_' + number

    def peek_friend(self):
        if len(self.friends) == self.num_dialog_friends:
            return random.choice(self.friends)
        while True:
            fr_id = random.choice(user_ids)
            if fr_id == self.user_id:
                return None
            self.friends.append(fr_id)
            return fr_id

    @task
    def dialog_send(self):
        friend = self.peek_friend()
        if friend is None:
            return

        with self.client.rename_request("/dialog/[friend]/send"):
            self.openapi_client.dialog_user_id_send_post(
                friend,
                dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    @task
    def dialog_list(self):
        friend = self.peek_friend()
        if friend is None:
            return

        with self.client.rename_request("/dialog/[friend]/list"):
            self.openapi_client.dialog_user_id_list_get(friend)


@events.init_command_line_parser.add_listener
def _(parser):
    parser.add_argument("--dialog-friends", type=int, env_var="DIALOG_FRIENDS",
                        dest="num_dialog_friends", default=2, help="Number of friends to dialog")


if __name__ == "__main__":
    run_single_user(DialogUser)
