from common.service import InvalidUserError, UserService

service = UserService()


def handler(event, context):
    try:
        # Synchronize confirmed Cognito user attributes into MongoDB
        user = service.create_from_cognito(event)
        print(f"Cognito user synchronized to MongoDB: {user.get('cognitoSub')}")

        # Return event payload unchanged as required by Cognito trigger
        return event

    except InvalidUserError as exc:
        print(f"Invalid Cognito user event: {exc}")
        raise

    except Exception as exc:
        print(f"Cognito user persistence error: {exc}")
        raise
