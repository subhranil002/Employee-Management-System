import json


def success(
    status_code,
    message,
    data=None,
):
    return {
        "success": True,
        "statusCode": status_code,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "message": message,
            "data": data,
        }),
    }


def error(
    status_code,
    message,
):
    return {
        "success": False,
        "statusCode": status_code,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "message": message,
            "data": None,
        }),
    }