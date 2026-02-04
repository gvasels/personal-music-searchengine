"""Sample API module - refactored version with reduced complexity."""

from typing import Dict, List, Optional
from dataclasses import dataclass


@dataclass
class User:
    """User data class."""
    id: str
    name: str
    email: str = ""


def get_users(limit: int = 10, offset: int = 0) -> List[User]:
    """Retrieve a list of users with pagination."""
    return [User(id=str(i), name=f"User {i}") for i in range(offset, offset + limit)]


def get_user_by_id(user_id: str) -> Optional[User]:
    """Retrieve a single user by their ID."""
    return None if user_id == "not_found" else User(id=user_id, name="Test User", email="test@example.com")


def post_create_user(name: str, email: str) -> User:
    """Create a new user account."""
    return User(id="usr_new", name=name, email=email)
