"""Sample API module for testing the starter project setup."""

from typing import Dict, List, Optional


def get_users(limit: int = 10, offset: int = 0) -> List[Dict]:
    """
    Retrieve a list of users with pagination.

    Args:
        limit: Maximum number of users to return
        offset: Number of users to skip

    Returns:
        List of user dictionaries
    """
    # Mock implementation
    return [{"id": i, "name": f"User {i}"} for i in range(offset, offset + limit)]


def get_user_by_id(user_id: str) -> Optional[Dict]:
    """
    Retrieve a single user by their ID.

    Args:
        user_id: The unique identifier of the user

    Returns:
        User dictionary or None if not found
    """
    # Mock implementation
    if user_id == "not_found":
        return None
    return {"id": user_id, "name": "Test User", "email": "test@example.com"}


def post_create_user(name: str, email: str) -> Dict:
    """
    Create a new user account.

    Args:
        name: User's display name
        email: User's email address

    Returns:
        Created user dictionary with generated ID
    """
    # Mock implementation
    return {"id": "usr_new", "name": name, "email": email}


def validate_input(data: Dict) -> bool:
    """Validate input data - not an API endpoint."""
    if not data:
        return False
    for key in ["name", "email"]:
        if key not in data:
            return False
    return True


class UserService:
    """Service class for user operations."""

    def __init__(self):
        self.users = {}

    def add_user(self, user: Dict) -> None:
        """Add a user to the service."""
        self.users[user["id"]] = user

    def get_user(self, user_id: str) -> Optional[Dict]:
        """Get a user from the service."""
        return self.users.get(user_id)
