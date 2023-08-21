# flake8: noqa

# import all models into this package
# if you have many models here with many references from one model to another this may
# raise a RecursionError
# to avoid this, import only the models that you directly need like:
# from openapi_client_dialog.model.pet import Pet
# or import this package, but before doing it, use:
# import sys
# sys.setrecursionlimit(n)

from openapi_client_dialog.model.dialog_message import DialogMessage
from openapi_client_dialog.model.dialog_user_id_send_post500_response import DialogUserIdSendPost500Response
from openapi_client_dialog.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
